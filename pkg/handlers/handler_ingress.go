package handlers

import (
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/ingress"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
)

func HandleCreateIngress(ingressManager *ingress.IngressManager, roomManager *room.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roomID := r.URL.Query().Get("roomID")

		room, exists := roomManager.GetRoom(roomID)
		if !exists {
			roomManager.CreateRoom(roomID)
			return
		}

		ingress, err := ingressManager.CreateIngress(room)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": map[string]string{
				"ingressId": ingress.IngressId,
				"url":       "rtmp://localhost:1935/live/",
				"streamKey": ingress.StreamKey,
			},
		})
	}
}

func HandleRemoveIngress(ingressManager *ingress.IngressManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingressId := r.URL.Query().Get("ingressId")

		err := ingressManager.RemoveIngress(ingressId)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
