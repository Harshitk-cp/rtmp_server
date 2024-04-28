package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
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

	go func() {
		// port := 1935 // Example RTMP port
		rtmpServer := rtmp.NewRTMPServer()
		err := rtmpServer.Start(&config.Config{}, func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error) {
			// Handle the onPublish callback here if needed
			return nil, nil, nil
		})
		if err != nil {
			log.Fatalf("Failed to start RTMP server: %v", err)
		}
	}()

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
