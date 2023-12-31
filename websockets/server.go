package websockets

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {

	// rooms holds all current rooms.
	Rooms map[*Room]bool

	// close is a channel for rooms wishing to close.
	Close chan *Room

	// newRoom is a channel for rooms wishing to be created.
	Room chan *Room
}

func NewServer() *Server {
	return &Server{
		Rooms: make(map[*Room]bool),
		Close: make(chan *Room),
		Room:  make(chan *Room),
	}
}

func (s *Server) GetRoom(roomId string) (map[*Client]bool, error) {
	for room := range s.Rooms {
		if room.id == roomId {
			return room.clients, nil
		}
	}
	return nil, nil

}

func (s *Server) Run() {
	for {
		select {
		case room := <-s.Room:
			s.Rooms[room] = true
		case room := <-s.Close:
			for cl := range room.clients {
				room.leave <- cl
			}
			s.Rooms[room] = false
		}
	}
}

func (s *Server) CreateRoom(roomId string) *Room {
	room := NewRoom(roomId)
	s.Room <- room
	return room
}

func (s *Server) JoinRoom(roomId string, client *Client) error {
	fmt.Println("in join room")
	fmt.Println(roomId)
	for room := range s.Rooms {
		// room exists so just join the client
		if room.id == roomId {
			fmt.Println("joining room")
			for cl := range room.clients {
				if cl.UserId == client.UserId {
					fmt.Println("user already exists")
					return fmt.Errorf("user already exists")
				}
			}
			client.Room = room
			room.join <- client
			return nil
		}
	}
	// if room doesn't exist create room and join the client
	room := s.CreateRoom(roomId)
	go room.run()
	client.Room = room
	room.join <- client
	fmt.Println(room)
	fmt.Println(s.Rooms)
	return nil

}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *Server) ServeHTTP(ctx *gin.Context, userId string) (*Client, error) {
	if userId == "undefined" {
		return nil, fmt.Errorf("user id is undefined")
	}
	// TODO : Check only trusted origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return nil, err
	}

	client := &Client{
		Socket:  socket,
		Receive: make(chan string, messageBufferSize),
		UserId:  userId,
	}
	go client.Read()
	go client.Write()

	// fmt.Println("New client: ", client, socket)

	return client, nil
	// fmt.Println(client)

	// r.join <- client
	// defer func() { r.leave <- client }()
	// go client.write()
	// client.read()
}
