package room

import (
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/utils"
)

type Room struct {
	ID           string
	Participants map[string]*Participant
	chatMessages chan ChatMessage
	streamKeys   map[string]bool
	mutex        sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:           id,
		Participants: make(map[string]*Participant),
		chatMessages: make(chan ChatMessage),
		streamKeys:   make(map[string]bool),
	}
}

func (r *Room) GenerateStreamKey() string {
	for {
		streamKey := utils.NewGuid("STREAM_KEY_")
		r.mutex.RLock()
		if _, exists := r.streamKeys[streamKey]; !exists {
			r.mutex.RUnlock()
			r.mutex.Lock()
			r.streamKeys[streamKey] = true
			r.mutex.Unlock()
			return streamKey
		}
		r.mutex.RUnlock()
	}
}

func (r *Room) CreateParticipant(id string) (*Participant, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Participants[id]; exists {
		return nil, errors.ErrParticipantExists
	}

	participant := &Participant{
		ID:     id,
		RoomID: r.ID,
	}

	r.Participants[id] = participant
	go participant.handleChatMessages(r.chatMessages)

	return participant, nil
}

func (r *Room) RemoveParticipant(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	participant, exists := r.Participants[id]
	if !exists {
		return errors.ErrParticipantNotFound
	}

	delete(r.Participants, id)
	participant.close()

	return nil
}

func (r *Room) SendChatMessage(message ChatMessage) {
	r.chatMessages <- message
}
