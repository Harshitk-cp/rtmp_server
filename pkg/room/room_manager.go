package room

import "sync"

type RoomManager struct {
	rooms map[string]*Room
	mutex sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) CreateRoom(id string) *Room {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room := NewRoom(id)
	rm.rooms[id] = room
	return room
}

func (rm *RoomManager) GetRoom(id string) (*Room, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	room, exists := rm.rooms[id]
	return room, exists
}

func (rm *RoomManager) RemoveRoom(id string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	delete(rm.rooms, id)
}

func (rm *RoomManager) JoinRoom(roomID string) (*Room, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		room = NewRoom(roomID)
		rm.rooms[roomID] = room
	}

	return room, nil
}
