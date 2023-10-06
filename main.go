package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	server "github.com/anudeep652/skribbl-clone-backend/websockets"
)

func generateRandomID() []byte {
	id, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func readFromBody(r *http.Request) string {
	resp, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(resp)
}

func main() {

	s := server.NewServer()
	go s.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client := s.ServeHTTP(w, r)
		fmt.Println(client)

	})
	http.HandleFunc("/new-room", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(readFromBody(r))
		// TODO : Read roomId from body and create a new room with that id, not with random id
		roomId := generateRandomID()
		s.NewRoom(roomId)
	})

	http.HandleFunc("/join-room", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(readFromBody(r))
		// TODO : Read roomId from body and join the room with that id

	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
