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

	// Get the room name and ingress ID from the query parameters
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

	// Add the new client to the room
	room.Clients[conn] = true

	// Handle incoming messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Broadcast the message to all clients in the room
		room.Broadcast(message)

		// Handle different message types if needed
		if messageType == websocket.CloseMessage {
			break
		}
	}

	// Remove the client from the room when disconnected
	delete(room.Clients, conn)
	conn.Close()
}
