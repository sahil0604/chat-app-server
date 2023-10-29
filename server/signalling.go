package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap

// CreateRoomRequestHandler create a room and return room id
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "* ")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
	log.Println(AllRooms.Map)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			client.Mutex.Lock()
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Fatal(err, "line 48 err")
					client.Conn.Close()
				}
			}
			client.Mutex.Unlock() // Unlock the mutex after the write operation
		}
	}
}

// JoinRoomRequestHandler will join the client in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg
		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("read Error:", err)
		}
		msg.Client = ws
		msg.RoomID = roomID[0]

		broadcast <- msg
	}
}
