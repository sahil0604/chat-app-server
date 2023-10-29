package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"video-chat-app/server"
)

func main() {
	server.AllRooms.Init()
	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	var port = envPortOr("3000")

	go server.Broadcaster()

	fmt.Println("Using port:", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}
