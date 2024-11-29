package controllers

import (
	"encoding/json"
	"html"
	"net/http"
	"sync"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GET_Threads(c *gin.Context) {
	threadList := make([]models.Thread, 0)

	// Grabbing classroom
	classroom, err := database.GetClassroom(c.Param("classroom_id"))
	if err != nil {
		utils.RouteError(c, "there was a problem")
		return
	}

	// Getting all the threads
	for _, threadId := range classroom.ThreadIDs {
		thread, err := database.GetThread(threadId)
		if err != nil {
			utils.RouteError(c, "there was a problem")
			return
		}

		threadList = append(threadList, *thread)
	}

	c.HTML(http.StatusOK, "class_board.html", gin.H{
		"IsLoggedIn": true,
		"Username":   c.GetString("username"),
		"ThreadList": threadList,
	})
}

// All open WebSockets, also a mutex to prevent race conditions
var openThreadsSockets = make(map[*websocket.Conn]string)
var openThreadsSocketsMutex sync.Mutex

// Broadcasts a message to all open sockets of a particular class ID
func broadcastThreadsSockets(data []byte, id string) {
	openThreadsSocketsMutex.Lock()
	defer openThreadsSocketsMutex.Unlock()

	for socket := range openThreadsSockets {
		// This open board isn't the same one we want to broadcast to
		if openThreadsSockets[socket] != id {
			continue
		}

		err := socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			socket.Close()
			delete(openThreadsSockets, socket)
		}
	}
}

func GET_ThreadsWebSocket(c *gin.Context) {
	// Upgrading connection to WebSocket
	socket, err := utils.WebSocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.RouteError(c, "there was a problem")
		return
	}
	defer socket.Close()

	// Add new WebSocket to open sockets
	openThreadsSocketsMutex.Lock()
	openThreadsSockets[socket] = c.Param("classroom_id")
	openThreadsSocketsMutex.Unlock()
	defer func() {
		openThreadsSocketsMutex.Lock()
		delete(openThreadsSockets, socket)
		openThreadsSocketsMutex.Unlock()
	}()

	errorCount := 0
	for {
		msgType, data, err := socket.ReadMessage()
		if err != nil {
			if errorCount++; errorCount == 20 {
				break
			}
			continue
		}

		if msgType == websocket.TextMessage {
			var message map[string]interface{}

			err = json.Unmarshal(data, &message)
			if err != nil {
				if errorCount++; errorCount == 20 {
					break
				}
				continue
			}

			switch message["Type"].(string) {
			case "newThread":
				if message["Title"].(string) == "" || message["Comment"].(string) == "" {
					break
				}

				thread, err := database.AddThread(c.GetString("username"), message["Title"].(string), c.Param("classroom_id"), message["Comment"].(string))
				if err != nil {
					break
				}

				// MongoDB handles escaping HTML, but since we are directly
				// sending the message back we need to manually escape
				message["Title"] = html.EscapeString(message["Title"].(string))

				message["ID"] = thread.ID
				message["Date"] = thread.Date
				message["AuthorImageURL"] = "/" + database.GetUserPFPPath(thread.AuthorID)
				message["AuthorID"] = thread.AuthorID
				message["OpenSockets"] = len(openThreadsSockets)

				jsonBytes, err := json.Marshal(message)
				if err != nil {
					break
				}

				broadcastThreadsSockets(jsonBytes, openThreadsSockets[socket])
				break

			default:
				break
			}
		} else {
			break
		}
	}
}
