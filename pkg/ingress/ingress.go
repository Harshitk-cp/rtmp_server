package ingress

import (
	"log"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/media"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
)

type Ingress struct {
	config       *config.Config
	rtmpServer   *rtmp.RTMPServer
	transcoder   *media.Transcoder
	webrtcSignal *media.WebRTCSignaling
	roomManager  *room.RoomManager
	streamKey    string
	resourceID   string
	room         *room.Room
	participant  *room.Participant
}

func NewIngress(config *config.Config, rtmpServer *rtmp.RTMPServer, streamKey string, r *room.Room, participant *room.Participant) *Ingress {
	return &Ingress{
		config:       config,
		rtmpServer:   rtmpServer,
		transcoder:   media.NewTranscoder(),
		webrtcSignal: media.NewWebRTCSignaling(),
		room:         r,
		participant:  participant,
	}
}

func (i *Ingress) SetParticipant(participant *room.Participant) {
	i.participant = participant
}

func (i *Ingress) Start() error {

	if i.transcoder == nil {
		i.transcoder = media.NewTranscoder()
	}

	go func() {
		err := i.transcoder.Start(i.rtmpServer, i.participant.ID)
		if err != nil {
			log.Print(err)
		}
	}()
	go func() {
		err := i.webrtcSignal.Start(i.participant)
		if err != nil {
			log.Print(err)
		}
	}()

	return nil
}

func (i *Ingress) Stop() {
	if i.transcoder != nil {
		i.transcoder.Stop()
	}
	if i.webrtcSignal != nil {
		i.webrtcSignal.Stop()
	}
	i.room.RemoveParticipant(i.participant.ID)
}
