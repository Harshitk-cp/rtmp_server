package room

import (
	"sync"

	protoutils "github.com/livekit/protocol/utils"
)

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

	ingressId := protoutils.NewGuid(protoutils.IngressPrefix)
	streamKey := protoutils.NewGuid(protoutils.RTMPResourcePrefix)

	room := NewRoom(id, streamKey, ingressId)
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
