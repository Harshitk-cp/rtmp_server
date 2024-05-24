package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/Harshitk-cp/rtmp_server/pkg/service"
	"github.com/joho/godotenv"
	"github.com/livekit/protocol/logger"
)

func main() {
	godotenv.Load()
	configFile := os.Getenv("CONFIG_FILE_PATH")

	conf, err := config.LoadFromFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to load Configuration: %s", err))
	}

	rtmpServer := rtmp.NewRTMPServer()
	relay := service.NewRelay(rtmpServer)

	err = rtmpServer.Start(conf, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed starting RTMP server %s", err))
	}
	err = relay.Start(conf)
	if err != nil {
		panic(fmt.Sprintf("Failed starting RTMP relay %s", err))
	}

	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, syscall.SIGINT)

	sig := <-killChan
	logger.Infow("exit requested, shutting down", "signal", sig)
	rtmpServer.Stop()
}