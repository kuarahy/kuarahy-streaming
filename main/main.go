package main

import (
	"log"
	"net/http"
	"twitch-notifier/internal/api"
	"twitch-notifier/internal/config"
	"twitch-notifier/internat/eventsub"
)

func main() {
	// Vite
	fs := http.FileServer(http.Dir("./frontend/dist")) // Vite

	// config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Services
	eventSub := eventsub.New(cfg.Twitch)
	obsManager := obsNewManager(cfg.OBS)
	apiServer := api.NewServer(cfg.Server, eventSub, obsManager)

	// Starting HTTP Server
	/* We need this for real-time calls to the API */
	go func() {
		log.Printf("Starting API server on %s", cfg.Server.Address)
		if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	// Start OBS Websocket
	/* This makes OBS connect to the local webserver */
	go func() {
		log.Printf("Starting OBS Webserver on %s", cfg.OBS.WSAddress)
		if err := obsManager.StartWSServer(); err != nil {
			log.Fatalf("OBS WebSocket error: %v", err)
		}
	}()

	// Calling Twitch EventSub service
	if err := eventSub.Connect(); err != nil {
		log.Fatalf("Failed to connect to Twitch: %v", err)
	}

	// Wait for suhtdown signal
	<-apiServer.Done()
}
