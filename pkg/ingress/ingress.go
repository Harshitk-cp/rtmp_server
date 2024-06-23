package ingress

import (
	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"

	"github.com/sirupsen/logrus"
)

type Ingress struct {
	IngressId   string
	StreamKey   string
	config      *config.Config
	sfuServer   *rtmp.SFUServer
	room        *room.Room
	participant *room.Participant
}

func NewIngress(config *config.Config, sfuServer *rtmp.SFUServer, streamKey string, r *room.Room, ingressId string) *Ingress {
	return &Ingress{
		IngressId: ingressId,
		StreamKey: streamKey,
		config:    config,
		sfuServer: sfuServer,
		room:      r,
	}
}

func (i *Ingress) Start() error {
	// go i.sfuServer.SendRTMPToWebRTC(i.room.ServerPeer)

	return nil
}

func (i *Ingress) Stop() {
	logrus.Error("ingress stopped")
	i.room.RemoveParticipant(i.participant.ID)
}
