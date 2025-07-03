package main

import (
	"log"
	"net/http"

	"./internal/api"
	"./internal/config"
	"./internal/eventsub"
)

func main() {
	// ===== Serve React static files go here =====
	// For Vite:
	fs := http.FileServer(http.Dir("./frontend/dist"))

	http.Handle("/", fs) // Serves index.html for all unmatched routes

	// ===== API routes go here =====
	http.HandleFunc("/api/events", handleEvents)
	http.HandleFunc("/api/config", handleConfig)

	// ===== (3) Start the server =====
	log.Println("Server starting on :5173...") // Go API is under 7777, don't get confused
	log.Fatal(http.ListenAndServe(":5173", nil))

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
