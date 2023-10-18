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
func RoomResponseJson(room map[*server.Client]bool) []RoomResponse {
	var resp []RoomResponse = []RoomResponse{}
	for client := range room {
		resp = append(resp, RoomResponse{UserId: client.UserId})
	}
	return resp

}

type RoomResponse struct {
	UserId string `json:"userId"`
}

func main() {

	s := server.NewServer()
	go s.Run()

	r := gin.Default()

	r.GET("/join-room/:roomId", func(ctx *gin.Context) {
		roomId := ctx.Params.ByName("roomId")
		userId := ctx.Query("user")
		fmt.Println(roomId)
		fmt.Println(userId)
		client, err := s.ServeHTTP(ctx, userId)
		if err != nil {
			return
		}
		fmt.Println("came here")

		err = s.JoinRoom(roomId, client)
		if err != nil {
			fmt.Println("exists")
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "user already exists",
			})
		}

		// room, err := s.GetRoom(roomId)
		// if err != nil {
		// 	ctx.JSON(http.StatusNotFound, gin.H{
		// 		"message": "Room not found",
		// 	})
		// }
		// fmt.Println("room:", room)
		// resp := RoomResponseJson(room)
		// fmt.Println("room response:", resp)
		// json, err := json.Marshal(resp)
		// if err != nil {
		// 	return
		// }
		client.Room.Forward <- ""
	})

	r.GET("/:roomId", func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		roomId := ctx.Params.ByName("roomId")
		room, err := s.GetRoom(roomId)
		var resp = RoomResponseJson(room)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Room not found",
			})
		}
		fmt.Println(room)
		ctx.JSON(http.StatusOK, gin.H{

			"users": resp,
		})

	})

	r.Run(":8000")
}
