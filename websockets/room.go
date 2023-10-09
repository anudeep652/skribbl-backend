package websockets

import "fmt"

type Room struct {

	// unique id for every room
	id string

	// clients holds all current clients in this room.
	clients map[*Client]bool

	// join is a channel for clients wishing to join the room.
	join chan *Client

	// leave is a channel for clients wishing to leave the room.
	leave chan *Client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.
	forward chan []byte
}

func NewRoom(id string) *Room {

	return &Room{
		id:      id,
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.Receive)
		case msg := <-r.forward:
			fmt.Print("in room run")
			for client := range r.clients {
				client.Receive <- msg
			}
		}
	}
}
