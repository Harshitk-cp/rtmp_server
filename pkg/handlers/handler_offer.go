package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/media"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/sirupsen/logrus"
)

type SDPExchange struct {
	SDP string `json:"sdp"`
}

func HandleOffer(w http.ResponseWriter, r *http.Request) {

	var sdpExchange SDPExchange
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &sdpExchange); err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusInternalServerError)
		return
	}
	logrus.Error(body)

	signaling := media.NewWebRTCSignaling()
	participant := &room.Participant{}
	answerSDP, err := signaling.Start(participant, sdpExchange.SDP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start WebRTC signaling: %v", err), http.StatusInternalServerError)
		return
	}

	response := SDPExchange{SDP: answerSDP}
	respondWithJSON(w, http.StatusOK, response)
}
