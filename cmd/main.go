package main

import (
	"net/http"
	"os"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/handlers"
	"github.com/Harshitk-cp/rtmp_server/pkg/ingress"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/Harshitk-cp/rtmp_server/pkg/service"
	"github.com/Harshitk-cp/rtmp_server/pkg/stats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-gst/go-gst/gst"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/checkHealth", handlers.HandleReadiness)
	v1Router.Get("/err", handlers.HandleErr)
	v1Router.Get("/ws/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		logrus.Print("roomId: ", roomID)
	})

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	conf, err := loadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	rtmpServer := rtmp.NewRTMPServer()
	relay := service.NewRelay()

	gst.Init(nil)

	go func() {
		err := rtmpServer.Start(1935, conf, func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error) {

			rm := room.NewRoomManager()
			r, _ := rm.GetRoom(streamKey)
			if r == nil {
				r = rm.CreateRoom(streamKey)
			}
			logrus.Printf("Stream key: %v", streamKey)

			participant, _ := r.CreateParticipant(resourceId)

			ingressManager := ingress.NewIngressManager(conf, rtmpServer)
			ingress, err := ingressManager.CreateIngress(streamKey, r, participant)
			if err != nil {
				logrus.Print(err.Error())
			}

			go ingress.Start()
			return nil, nil, nil
		})
		if err != nil {
			logrus.Fatalf("Failed to start RTMP server: %v", err)
		}
	}()

	go func() {
		err := relay.Start(conf, 9090)
		if err != nil {
			logrus.Fatalf("Failed to start RTMP relay: %v", err)
		}
	}()

	logrus.Printf("Server starting on port %v", portString)

	servErr := srv.ListenAndServe()
	if servErr != nil {
		logrus.Fatal(servErr)
	}

}

func loadConfig() (*config.Config, error) {
	configFile := os.Getenv("CONFIG_FILE_PATH")
	if configFile == "" {
		return nil, errors.ErrNoConfig
	}

	conf, err := config.LoadFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
