package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Room struct {
	Name      string
	IngressID string
	Clients   map[*websocket.Conn]bool
}

var Rooms = make(map[string]*Room)

func CreateRoom(name, ingressID string) *Room {
	room := &Room{
		Name:      name,
		IngressID: ingressID,
		Clients:   make(map[*websocket.Conn]bool),
	}
	key := name + "::" + ingressID
	Rooms[key] = room
	return room
}

func (r *Room) Broadcast(message []byte) {
	for client := range r.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(r.Clients, client)
		}
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomName := r.URL.Query().Get("room")
	ingressID := r.URL.Query().Get("ingressID")

	if roomName == "" || ingressID == "" {
		// Handle error: room name or ingress ID not provided
		conn.Close()
		return
	}

	key := roomName + "::" + ingressID
	room, ok := Rooms[key]
	if !ok {
		// Create a new room if it doesn't exist
		room = CreateRoom(roomName, ingressID)
	}

	room.Clients[conn] = true
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		room.Broadcast(message)

		if messageType == websocket.CloseMessage {
			break
		}
	}
	delete(room.Clients, conn)
	conn.Close()
}
