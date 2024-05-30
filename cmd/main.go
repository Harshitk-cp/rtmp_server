package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-gst/go-gst/gst"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/handlers"
	"github.com/Harshitk-cp/rtmp_server/pkg/ingress"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/Harshitk-cp/rtmp_server/pkg/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	router := setupRouter()

	conf, err := loadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}
	port := conf.HttpPort

	gst.Init(nil)

	rtmpServer := rtmp.NewRTMPServer()
	relay := service.NewRelay(rtmpServer, router)

	go startRTMPServer(rtmpServer, conf)

	err = relay.Start()
	if err != nil {
		logrus.Fatalf("Failed to start relay server: %v", err)
	}

	servErr := startHTTPServer(router, fmt.Sprintf(":%d", port))
	if servErr != nil {
		logrus.Fatal(servErr)
	}
}

func setupRouter() *chi.Mux {
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
	v1Router.Get("/health", handlers.HandleReadiness)
	v1Router.Get("/error", handlers.HandleErr)
	v1Router.Get("/ws/{roomID}", handlers.HandleWebSocket)

	router.Mount("/v1", v1Router)
	return router
}

func startHTTPServer(router *chi.Mux, port string) error {
	srv := &http.Server{
		Handler: router,
		Addr:    port,
	}
	logrus.Infof("Server starting on port %v", port)

	return srv.ListenAndServe()
}

func startRTMPServer(rtmpServer *rtmp.RTMPServer, conf *config.Config) {
	err := rtmpServer.Start(conf, func(streamKey, resourceId string) (*params.Params, error) {
		rm := room.NewRoomManager()
		r, exists := rm.GetRoom(streamKey)
		if !exists {
			r = rm.CreateRoom(streamKey)
		}
		logrus.Infof("Stream key: %v", streamKey)

		participant, err := r.CreateParticipant(resourceId)
		if err != nil {
			logrus.Errorf("Failed to create participant: %v", err)
			return nil, err
		}

		ingressManager := ingress.NewIngressManager(conf, rtmpServer)
		ingress, err := ingressManager.CreateIngress(streamKey, r, participant)
		if err != nil {
			logrus.Errorf("Failed to create ingress: %v", err)
			return nil, err
		}
		go ingress.Start()

		return nil, nil
	})

	if err != nil {
		logrus.Fatalf("Failed to start RTMP server: %v", err)
	}
}

func loadConfig() (*config.Config, error) {
	configFile := os.Getenv("CONFIG_FILE_PATH")
	conf, err := config.LoadFromFile(configFile)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
