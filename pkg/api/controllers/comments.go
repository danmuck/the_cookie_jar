package controllers

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"sync"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GET_Comments(c *gin.Context) {
	commentList := make([]models.Comment, 0)

	// Grabbing thread
	thread, err := database.GetThread(c.Param("thread_id"))
	if err != nil {
		utils.RouteError(c, "this content does not exist")
		return
	}

	// Getting all the comments
	for _, commentId := range thread.CommentIDs {
		comment, err := database.GetComment(commentId)
		if err != nil {
			utils.RouteError(c, commentId)
			return
		}

		commentList = append(commentList, *comment)
	}

	c.HTML(http.StatusOK, "class_board_thread.html", gin.H{
		"IsLoggedIn":       true,
		"Username":         c.GetString("username"),
		"IsClassProfessor": c.GetBool("isClassProfessor"),
		"ThreadTitle":      thread.Title,
		"CommentList":      commentList,
	})
}

// All open WebSockets, also a mutex to prevent race conditions
var openCommentsSockets = make(map[*websocket.Conn]string)
var openCommentsSocketsMutex sync.Mutex

// Broadcasts a message to all open sockets of a particular thread ID
func broadcastCommentsSocket(data []byte, id string) {
	openCommentsSocketsMutex.Lock()
	defer openCommentsSocketsMutex.Unlock()

	for socket := range openCommentsSockets {
		// This open thread isn't the same one we want to broadcast to
		if openCommentsSockets[socket] != id {
			continue
		}

		err := socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			socket.Close()
			delete(openCommentsSockets, socket)
		}
	}
}

func GET_CommentsWebSocket(c *gin.Context) {
	// Upgrading connection to WebSocket
	socket, err := utils.WebSocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.RouteError(c, "there was a problem")
		return
	}
	defer socket.Close()

	// Add new WebSocket to open sockets
	openCommentsSocketsMutex.Lock()
	openCommentsSockets[socket] = c.Param("thread_id")
	openCommentsSocketsMutex.Unlock()
	defer func() {
		openCommentsSocketsMutex.Lock()
		delete(openCommentsSockets, socket)
		openCommentsSocketsMutex.Unlock()
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
			case "newComment":
				if message["Text"].(string) == "" {
					break
				}

				comment, err := database.AddComment(c.GetString("username"), message["Text"].(string), c.Param("thread_id"))
				if err != nil {
					break
				}

				// MongoDB handles escaping HTML, but since we are directly
				// sending the message back we need to manually escape
				message["Text"] = html.EscapeString(message["Text"].(string))

				message["ID"] = comment.ID
				message["AuthorImageURL"] = "/" + database.GetUserPFPPath(comment.AuthorID)
				message["AuthorID"] = comment.AuthorID
				message["OpenSockets"] = len(openCommentsSockets)

				jsonBytes, err := json.Marshal(message)
				if err != nil {
					break
				}

				broadcastCommentsSocket(jsonBytes, openCommentsSockets[socket])

			case "likeComment":
				if message["ID"].(string) == "" {
					break
				}

				userLiked, err := database.IsUserLiked(message["ID"].(string), c.GetString("username"))
				if err != nil {
					break
				}

				message["Liked"] = strconv.FormatBool(userLiked)
				jsonBytes, err := json.Marshal(message)
				if err != nil {
					break
				}

				broadcastCommentsSocket(jsonBytes, openCommentsSockets[socket])

			case "editComment":
				if message["ID"].(string) == "" || message["Text"].(string) == "" {
					break
				}

				// Making sure user isn't trying to edit other people's
				// comments unless they are the professor
				err = database.IsUsersComment(message["ID"].(string), c.GetString("username"))
				if !c.GetBool("isClassProfessor") && err != nil {
					break
				}

				// Updating their comment in database
				err = database.UpdateCommentText(message["ID"].(string), message["Text"].(string))
				if err != nil {
					break
				}

				jsonBytes, err := json.Marshal(message)
				if !c.GetBool("isClassProfessor") && err != nil {
					break
				}

				broadcastCommentsSocket(jsonBytes, openCommentsSockets[socket])
			}
		} else {
			break
		}
	}
}
