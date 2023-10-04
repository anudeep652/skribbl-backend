package main

import (
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)
	go server.Serve()
	defer server.Close()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	http.ListenAndServe(":8000", nil)
}
