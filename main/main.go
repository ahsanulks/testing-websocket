package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

type MessagePayload struct {
	Username string `json:"Username"`
	Type     string `json:"Type"`
	Message  string `json:"Message"`
}

type Conn struct {
	Connection *websocket.Conn
	Username   string
}

var connections []*Conn

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) {
		data := make(map[string]interface{})
		c.Render("./view/index.html", data)
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		mt, msg, _ := c.ReadMessage()
		var message MessagePayload
		json.Unmarshal(msg, &message)
		currentConnection := Conn{
			Connection: c,
			Username:   message.Username,
		}
		connections = append(connections, &currentConnection)
		messageSend := MessagePayload{
			Username: message.Username,
			Type:     "New",
			Message:  "joining the chat",
		}
		jsonMessage, _ := json.Marshal(messageSend)
		broardcast(&currentConnection, mt, jsonMessage)
		handleConnection(&currentConnection)
	}))

	app.Listen(3000)
}

func handleConnection(currentConnection *Conn) {
	for {
		mt, msg, err := currentConnection.Connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			removeConnection(currentConnection)
			messageSend := MessagePayload{
				Username: currentConnection.Username,
				Type:     "Leave",
				Message:  "Leave the chat",
			}
			jsonMessage, _ := json.Marshal(messageSend)
			broardcast(currentConnection, 1, jsonMessage)
			break
		}
		var a MessagePayload
		json.Unmarshal(msg, &a)
		if a.Type == "message" {
			a.Type = "message"
			jsonMessage, _ := json.Marshal(a)
			broardcast(currentConnection, mt, jsonMessage)
		} else {
			jsonMessage, _ := json.Marshal(a.Username + " leave the chat")
			broardcast(currentConnection, mt, jsonMessage)
		}
	}
}

func broardcast(currentConnection *Conn, messageType int, message []byte) {
	connects := connections
	for _, v := range connects {
		if v == currentConnection {
			continue
		}
		v.Connection.WriteMessage(messageType, message)
	}
}

func removeConnection(currentConnection *Conn) {
	var filtered []*Conn
	for _, v := range connections {
		if v != currentConnection {
			filtered = append(filtered, v)
		}
	}
	connections = filtered
}
