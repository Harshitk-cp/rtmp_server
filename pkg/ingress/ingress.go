package ingress

import (
	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"

	"github.com/sirupsen/logrus"
)

type Ingress struct {
	config      *config.Config
	sfuServer   *rtmp.SFUServer
	roomManager *room.RoomManager
	streamKey   string
	resourceID  string
	room        *room.Room
	participant *room.Participant
}

func NewIngress(config *config.Config, sfuServer *rtmp.SFUServer, streamKey string, r *room.Room) *Ingress {
	return &Ingress{
		config:    config,
		sfuServer: sfuServer,
		room:      r,
	}
}

func (i *Ingress) SetParticipant(participant *room.Participant) {
	i.participant = participant
}

func (i *Ingress) Start() error {
	return nil
}

func (i *Ingress) Stop() {
	logrus.Error("ingress stopped")
	i.room.RemoveParticipant(i.participant.ID)
}
