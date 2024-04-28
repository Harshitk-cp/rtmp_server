package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/handlers"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/Harshitk-cp/rtmp_server/pkg/stats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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
	v1Router.Get("/es", handlers.HandleWebSocket)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	go func() {
		rtmpServer := rtmp.NewRTMPServer()
		err := rtmpServer.Start(conf, func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error) {
			return nil, nil, nil
		})
		if err != nil {
			log.Fatalf("Failed to start RTMP server: %v", err)
		}
	}()

	log.Printf("Server starting on port %v", portString)

	servErr := srv.ListenAndServe()
	if servErr != nil {
		log.Fatal(servErr)
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
