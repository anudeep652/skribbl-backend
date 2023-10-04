package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns  map[*websocket.Conn]bool
	length int
}

func (s *Server) connection(ws *websocket.Conn) {
	s.length++
	s.conns[ws] = true
	userMsg := fmt.Sprintf("New user joined: #%d", s.length)
	fmt.Println(userMsg)

	s.BroadCastMsg([]byte(userMsg))
	s.ReadMsg(ws)
}

func (s *Server) ReadMsg(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		msg := buf[:n]

		fmt.Println("Received message from", ws.RemoteAddr(), string(msg))
		s.BroadCastMsg(msg)
	}

}

func (s *Server) BroadCastMsg(msg []byte) {
	for ws := range s.conns {

		go func(ws *websocket.Conn) {

			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error sending message to", ws.RemoteAddr(), err.Error())
				delete(s.conns, ws)
				return
			}
		}(ws)
	}
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.connection))
	http.ListenAndServe(":8000", nil)

}
