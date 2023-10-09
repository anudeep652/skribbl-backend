package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	server "github.com/anudeep652/skribbl-clone-backend/websockets"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/join-room/:roomId", func(ctx *gin.Context) {
		roomId := ctx.Params.ByName("roomId")
		fmt.Println(roomId)
		client := s.ServeHTTP(ctx)

		s.JoinRoom(roomId, client)
		ctx.JSON(http.StatusOK, gin.H{})
	})

	r.Run(":8000")
}
