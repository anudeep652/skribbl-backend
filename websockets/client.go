package websockets

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type Client struct {
	UserId string `json:"userId"`

	// socket is the web socket for this client.
	Socket *websocket.Conn

	// receive is a channel to receive messages from other clients.
	Receive chan string

	// room is the room this client is chatting in.
	Room *Room
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(msg))
		c.Room.Forward <- string(msg)
	}
}

func (c *Client) Write() {
	fmt.Println("in write")
	defer c.Socket.Close()
	for msg := range c.Receive {
		fmt.Println("Sending message to client")

		err := c.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return
		}
	}
}
