package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {

	// rooms holds all current rooms.
	rooms map[*Room]bool

	// close is a channel for rooms wishing to close.
	close chan *Room

	// newRoom is a channel for rooms wishing to be created.
	newRoom chan *Room
}

func NewServer() *Server {
	return &Server{
		rooms:   make(map[*Room]bool),
		close:   make(chan *Room),
		newRoom: make(chan *Room),
	}
}

func (s *Server) Run() {
	for {
		select {
		case room := <-s.newRoom:
			s.rooms[room] = true
			go room.run()
		case room := <-s.close:
			for cl := range room.clients {
				room.leave <- cl
			}
			s.rooms[room] = false
		}
	}
}

func (s *Server) NewRoom(roomId []byte) {
	room := newRoom(roomId)
	s.newRoom <- room
}

func (s *Server) JoinRoom(roomId []byte, client *Client) {
	for room := range s.rooms {
		if string(room.id) == string(roomId) {
			room.join <- client
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) *Client {
	// TODO : Check only trusted origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return nil
	}

	client := &Client{
		Socket:  socket,
		Receive: make(chan []byte, messageBufferSize),
	}
	go client.Read()

	// fmt.Println("New client: ", client, socket)

	return client
	// fmt.Println(client)

	// r.join <- client
	// defer func() { r.leave <- client }()
	// go client.write()
	// client.read()
}
