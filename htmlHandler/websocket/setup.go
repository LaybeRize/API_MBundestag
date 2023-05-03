package websocket

import (
	"API_MBundestag/database_old"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Message represents a message in the chat room
type Message struct {
	Text   string `json:"Text"`
	SendBy string `json:"SendBy"`
}

// Room represents a chat room for two users
type Room struct {
	RoomToken        string
	User1Token       string
	UserName1        string
	User1LastMessage string
	User2Token       string
	UserName2        string
	User2LastMessage string
	User1Conn        *websocket.Conn
	User2Conn        *websocket.Conn
}

var rooms = map[string]*Room{}

func GetWebsocket(c *gin.Context) {
	token := c.Param("token")
	chat := database.Chat{}
	userToken := c.Param("user")
	err := chat.GetByID(token, userToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	room := getRoom(token)
	var b bool
	if b, room = setUpUser(room, conn, chat, userToken); b {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		var messageStruct Message
		if err := json.Unmarshal(message, &messageStruct); err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
		}

		room.sendMessagesToUser(room.User1Conn, messageStruct.Text, userToken)
		room.sendMessagesToUser(room.User2Conn, messageStruct.Text, userToken)
	}
	iterateMapAndCheckConnected()
}

func getRoom(token string) *Room {
	return rooms[token]
}

func (m *Room) sendMessagesToUser(conn *websocket.Conn, message string, userToken string) {
	if conn == nil {
		return
	}
	userName := m.UserName1
	if m.User1Token != userToken {
		m.User2LastMessage = message
		userName = m.UserName2
	} else {
		m.User1LastMessage = message
	}
	conn.WriteJSON(Message{
		Text:   message,
		SendBy: userName,
	})
}

func iterateMapAndCheckConnected() {
	for key, val := range rooms {
		if val.checkIfConnected() {
			delete(rooms, key)
		}
	}
}

func (m *Room) checkIfConnected() bool {
	if m.User1Conn != nil && m.User1Conn.WriteMessage(websocket.PingMessage, nil) != nil {
		m.User1Conn = nil
		m.UserName1 = ""
	}
	if m.User2Conn != nil && m.User2Conn.WriteMessage(websocket.PingMessage, nil) != nil {
		m.User2Conn = nil
		m.UserName2 = ""
	}
	if m.User1Conn == nil && m.User2Conn == nil {
		return true
	}
	return false
}

func setUpUser(m *Room, conn *websocket.Conn, chat database.Chat, myToken string) (bool, *Room) {
	myUsername, myLastMessage, otherLastMessage := getUserByToken(chat, myToken)
	if m == nil {
		m = &Room{RoomToken: chat.UUID,
			User1Token:       myToken,
			User1Conn:        conn,
			UserName1:        myUsername,
			User1LastMessage: myLastMessage,
			User2LastMessage: otherLastMessage,
		}
		rooms[chat.UUID] = m
		m.sendMessagesToUser(conn, m.User2LastMessage, "")
		m.sendMessagesToUser(conn, m.User1LastMessage, m.User1Token)
		return false, m
	}
	if m.UserName1 == myUsername || m.UserName2 == myUsername {
		return true, m
	}
	if m.UserName1 == "" {
		m.User1Token = myToken
		m.UserName1 = myUsername
		m.User1Conn = conn
		m.sendMessagesToUser(conn, m.User2LastMessage, m.User2Token)
		m.sendMessagesToUser(conn, m.User1LastMessage, m.User1Token)
		return false, m
	}
	m.User2Token = myToken
	m.UserName2 = myUsername
	m.User2Conn = conn
	m.sendMessagesToUser(conn, m.User1LastMessage, m.User1Token)
	m.sendMessagesToUser(conn, m.User2LastMessage, m.User2Token)
	return false, m
}

func getUserByToken(chat database.Chat, token string) (string, string, string) {
	if chat.User1 == token {
		return chat.Username1, chat.Message1, chat.Message2
	} else {
		return chat.Username2, chat.Message2, chat.Message1
	}
}
