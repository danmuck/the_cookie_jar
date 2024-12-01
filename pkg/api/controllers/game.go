package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Question struct {
	Text string

	CorrectOption int
	Options       []string
}
type Game struct {
	ClassID       string
	IsStarted     bool
	IsTimerHit    bool
	PlayerState   map[string]int
	PlayerScore   map[string]int
	PlayerSockets map[string]*websocket.Conn

	CurrentQuestion int
	Questions       []Question
	SecondsLeft     int

	Mutex sync.Mutex
}

var defaultQuestions = []Question{
	{
		Text:          "What is the capital of New York state?",
		CorrectOption: 2,
		Options:       []string{"Houston", "Washington DC", "Albany", "Buffalo"},
	},
	{
		Text:          "Which animal is known as the 'King of the Jungle'?",
		CorrectOption: 0,
		Options:       []string{"Lion", "Elephant", "Tiger", "Giraffe"},
	},
	{
		Text:          "Which planet is known as the Red Planet?",
		CorrectOption: 0,
		Options:       []string{"Mars", "Venus", "Jupiter", "Saturn"},
	},
	{
		Text:          "What is the tallest mountain in the world?",
		CorrectOption: 3,
		Options:       []string{"K2", "Mount Kilimanjaro", "Mount Fuji", "Mount Everest"},
	},
	{
		Text:          "What color is the 'Exclamation Point' in the Nintendo game Super Mario Bros.?",
		CorrectOption: 2,
		Options:       []string{"Green", "Red", "Yellow", "Blue"},
	},
	{
		Text:          "Which country is famous for creating the pizza?",
		CorrectOption: 1,
		Options:       []string{"France", "Italy", "Germany", "USA"},
	},
	{
		Text:          "How many colors are there in a rainbow?",
		CorrectOption: 4,
		Options:       []string{"6", "7", "8", "9"},
	},
	{
		Text:          "Which fruit is also known as a 'love apple'?",
		CorrectOption: 1,
		Options:       []string{"Orange", "Tomato", "Apple", "Banana"},
	},
	{
		Text:          "In which country can you find the Eiffel Tower?",
		CorrectOption: 0,
		Options:       []string{"France", "Italy", "England", "Spain"},
	},
	{
		Text:          "What was the first video game ever created?",
		CorrectOption: 2,
		Options:       []string{"Pac-Man", "Tetris", "Pong", "Space Invaders"},
	},
	{
		Text:          "How many legs does an octopus have?",
		CorrectOption: 0,
		Options:       []string{"8", "6", "10", "12"},
	},
	{
		Text:          "What is the only mammal capable of true flight?",
		CorrectOption: 3,
		Options:       []string{"Bat", "Swan", "Flying Squirrel", "Bird"},
	},
	{
		Text:          "What is the capital of Canada?",
		CorrectOption: 1,
		Options:       []string{"Toronto", "Ottawa", "Vancouver", "Montreal"},
	},
	{
		Text:          "Which of the following is NOT a primary color?",
		CorrectOption: 2,
		Options:       []string{"Red", "Blue", "Green", "Yellow"},
	},
	{
		Text:          "Which vegetable is known for making people cry when chopped?",
		CorrectOption: 0,
		Options:       []string{"Onion", "Garlic", "Carrot", "Cucumber"},
	},
}

var activeGames = make(map[string]*Game)

func (game *Game) addPlayer(player string, socket *websocket.Conn) {
	game.PlayerState[player] = 0
	game.PlayerScore[player] = 0
	game.PlayerSockets[player] = socket
}

func (game *Game) removePlayer(player string) {
	delete(game.PlayerState, player)
	delete(game.PlayerScore, player)
	delete(game.PlayerSockets, player)
}

func (game *Game) readyPlayer(player string) {
	state, exists := game.PlayerState[player]
	if !exists {
		game.PlayerState[player] = 1
		return
	}

	// If already ready un-ready them
	if state == 0 {
		game.PlayerState[player] = 1
	} else {
		game.PlayerState[player] = 0
	}
}

func (game *Game) checkAnswer(player string, answer int) {
	question := game.Questions[game.CurrentQuestion]
	if answer == question.CorrectOption {
		game.PlayerScore[player] += game.SecondsLeft*10 + 10
	}
}

func (game *Game) sendPlayer(player string, message map[string]interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	socket := game.PlayerSockets[player]
	err = socket.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		socket.Close()
		game.removePlayer(player)
		delete(openGamesSockets, socket)
	}
}

func (game *Game) sendAllPlayers(message map[string]interface{}) {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	broadcastGamesSockets(jsonBytes, game.ClassID)
}

func (game *Game) receive(c *gin.Context, socket *websocket.Conn, message map[string]interface{}) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	switch message["Type"].(string) {
	case "join":
		game.addPlayer(c.GetString("username"), socket)
		game.sendAllPlayers(map[string]interface{}{
			"Type":    "playerList",
			"Players": game.PlayerState,
		})

	case "start":
		if !c.GetBool("isClassProfessor") {
			game.sendAllPlayers(map[string]interface{}{
				"Type": "error",
				"Text": "Only the professor can start the game.",
			})
			break
		}
		game.IsStarted = true

	case "ready":
		game.readyPlayer(c.GetString("username"))
		game.sendAllPlayers(map[string]interface{}{
			"Type":    "playerList",
			"Players": game.PlayerState,
		})

	case "questionAnswer":
		answerStr, exists := message["Option"]
		if !exists {
			break
		}

		answerInt, err := strconv.Atoi(answerStr.(string))
		if err != nil {
			break
		}
		game.checkAnswer(c.GetString("username"), answerInt)
	}
}

func getGame(id string) *Game {
	game, exists := activeGames[id]
	if !exists {
		game = &Game{
			ClassID:       id,
			IsStarted:     false,
			IsTimerHit:    false,
			PlayerState:   make(map[string]int),
			PlayerScore:   make(map[string]int),
			PlayerSockets: make(map[string]*websocket.Conn),

			CurrentQuestion: 0,
			Questions:       defaultQuestions,
		}

		activeGames[id] = game
		go game.run()
	}

	return game
}

func (game *Game) gameCountdown(seconds int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for remaining := seconds - 1; remaining >= 0; remaining-- {
		game.SecondsLeft = remaining
		if game.IsStarted {
			game.sendAllPlayers(map[string]interface{}{
				"Type":    "questionCountdown",
				"Seconds": remaining,
			})
		} else {
			game.sendAllPlayers(map[string]interface{}{
				"Type":    "startCountdown",
				"Seconds": remaining,
			})
		}

		<-ticker.C
	}

	game.IsTimerHit = true
}

func (game *Game) run() {
	defer delete(activeGames, game.ClassID)

	for !game.IsStarted {
		time.Sleep(1 * time.Second)

		// If all players left just stop the game
		if len(game.PlayerState) == 0 {
			delete(activeGames, game.ClassID)
			return
		}
	}

	game.IsStarted = false
	game.gameCountdown(5)
	game.IsStarted = true

	for i, question := range game.Questions {
		game.CurrentQuestion = i

		// Sending the prep question (the title)
		game.sendAllPlayers(map[string]interface{}{
			"Type": "questionPrep",
			"Text": question.Text,
		})
		game.gameCountdown(3)

		// Then send the actual question
		game.sendAllPlayers(map[string]interface{}{
			"Type":    "question",
			"Text":    question.Text,
			"Options": question.Options,
		})
		game.gameCountdown(6)

		// Display leaderboard, points, and right answer
		game.sendAllPlayers(map[string]interface{}{
			"Type":   "leaderboard",
			"Scores": game.PlayerScore,
		})
		for player := range game.PlayerSockets {
			game.sendPlayer(player, map[string]interface{}{
				"Type":          "score",
				"Score":         game.PlayerScore[player],
				"CorrectOption": game.Questions[game.CurrentQuestion].CorrectOption,
			})
		}
		game.gameCountdown(6)
	}
}

func GET_Game(c *gin.Context) {
	c.HTML(http.StatusOK, "game.html", gin.H{
		"IsLoggedIn":       true,
		"Username":         c.GetString("username"),
		"IsClassProfessor": c.GetBool("isClassProfessor"),
	})
}

// All open WebSockets, also a mutex to prevent race conditions
var openGamesSockets = make(map[*websocket.Conn]string)
var openGamesSocketsMutex sync.Mutex

// Broadcasts a message to all open sockets of a particular class ID
func broadcastGamesSockets(data []byte, id string) {
	openGamesSocketsMutex.Lock()
	defer openGamesSocketsMutex.Unlock()

	for socket := range openGamesSockets {
		// This open board isn't the same one we want to broadcast to
		if openGamesSockets[socket] != id {
			continue
		}

		err := socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			socket.Close()
			delete(openGamesSockets, socket)
		}
	}
}

func GET_GameWebSocket(c *gin.Context) {
	// Upgrading connection to WebSocket
	socket, err := utils.WebSocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.RouteError(c, "there was a problem")
		return
	}
	defer socket.Close()

	// Add new WebSocket to open sockets
	openGamesSocketsMutex.Lock()
	openGamesSockets[socket] = c.Param("classroom_id")
	openGamesSocketsMutex.Unlock()
	defer func() {
		openGamesSocketsMutex.Lock()
		delete(openGamesSockets, socket)
		openGamesSocketsMutex.Unlock()
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

			game := getGame(openGamesSockets[socket])
			game.receive(c, socket, message)
		} else {
			break
		}
	}

	// If the game hasn't started and the connection is closed, remove the user
	game := getGame(openGamesSockets[socket])
	if !game.IsStarted {
		game.removePlayer(c.GetString("username"))
		game.sendAllPlayers(map[string]interface{}{
			"Type":    "playerList",
			"Players": game.PlayerState,
		})
	}
}
