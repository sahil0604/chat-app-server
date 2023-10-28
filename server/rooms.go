package server

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// participant describes a single entity in the hashMap
type Participant struct {
	Host bool
	Conn *websocket.Conn
}

// RoomMap is the main hashMao [roomID string] -> participant
type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

// Init initializes the Room Map  struct
func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

// Get will return array of participants in the room
func (r *RoomMap) Get(roomId string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomId]
}

// CreateRoom create a unique room id and return it and insert it into hashMap
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	roomID := string(b)

	r.Map[roomID] = []Participant{}
	return roomID
}

// InsertIntoRoom insert participant to the hash map
func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Println("Inserting into room with roomId", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

// DeleteRoom delete the room from hash map
func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}
