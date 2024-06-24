package handlers

import (
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
)

func HandleCreateIngress(roomManager *room.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roomID := r.URL.Query().Get("roomId")
		room, exists := roomManager.GetRoom(roomID)
		if exists {
			respondWithJSON(w, http.StatusOK, map[string]interface{}{
				"data": map[string]string{
					"ingressId": room.IngressID,
					"url":       "rtmp://localhost:1935/live/",
					"streamKey": room.StreamKey,
				},
			})
			return
		}
		newRoom := roomManager.CreateRoom(roomID)

		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": map[string]string{
				"ingressId": newRoom.IngressID,
				"url":       "rtmp://localhost:1935/live/",
				"streamKey": newRoom.StreamKey,
			},
		})
	}
}
