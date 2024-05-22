package room

import "sync"

type Participant struct {
	ID         string
	RoomID     string
	chatStream chan ChatMessage
	closed     bool
	mutex      sync.RWMutex
}

func (p *Participant) handleChatMessages(messages chan ChatMessage) {
	for msg := range messages {
		p.mutex.RLock()
		if p.closed {
			p.mutex.RUnlock()
			return
		}
		p.chatStream <- msg
		p.mutex.RUnlock()
	}
}

func (p *Participant) ReceiveChatMessage() <-chan ChatMessage {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.chatStream == nil {
		p.chatStream = make(chan ChatMessage)
	}

	return p.chatStream
}

func (p *Participant) close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.closed {
		p.closed = true
		close(p.chatStream)
	}
}
