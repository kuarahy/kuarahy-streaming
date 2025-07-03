// internal/eventsub/eventsub.go
package eventsub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nicklaw5/helix/v2"
)

type EventSub struct {
	client       *helix.Client
	callbackURL  string
	secret       string
	eventHandler func(eventType string, data interface{})
}

func New(cfg config.TwitchConfig) *EventSub {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &EventSub{
		client:      client,
		callbackURL: cfg.CallbackURL,
		secret:      cfg.WebhookSecret,
	}
}

func (es *EventSub) Connect() error {
	// Get or refresh OAuth token
	resp, err := es.client.RequestAppAccessToken([]string{})
	if err != nil {
		return fmt.Errorf("failed to get app access token: %v", err)
	}
	es.client.SetAppAccessToken(resp.Data.AccessToken)

	// Subscribe to events
	subscriptions := []struct {
		Type      string
		Version   string
		Condition map[string]string
	}{
		{Type: "channel.follow", Version: "2", Condition: map[string]string{"broadcaster_user_id": cfg.BroadcasterID}},
		{Type: "channel.subscribe", Version: "1", Condition: map[string]string{"broadcaster_user_id": cfg.BroadcasterID}},
		// Add more events as needed
	}

	for _, sub := range subscriptions {
		_, err := es.client.CreateEventSubSubscription(&helix.EventSubSubscription{
			Type:      sub.Type,
			Version:   sub.Version,
			Condition: sub.Condition,
			Transport: helix.EventSubTransport{
				Method:   "webhook",
				Callback: es.callbackURL,
				Secret:   es.secret,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to subscribe to %s: %v", sub.Type, err)
		}
	}

	return nil
}

func (es *EventSub) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify the webhook signature
	if !helix.VerifyEventSubNotification(es.secret, r) {
		http.Error(w, "Invalid signature", http.StatusForbidden)
		return
	}

	// Handle challenge response
	if r.Method == http.MethodGet {
		var challenge struct {
			Challenge string `json:"challenge"`
		}
		if err := json.NewDecoder(r.Body).Decode(&challenge); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(challenge.Challenge))
		return
	}

	// Process the event
	var notification helix.EventSubNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle different event types
	switch notification.Subscription.Type {
	case "channel.follow":
		var event helix.EventSubChannelFollowEvent
		if err := json.Unmarshal(notification.Event, &event); err != nil {
			log.Printf("Error unmarshaling follow event: %v", err)
			return
		}
		es.eventHandler("follow", event)
	case "channel.subscribe":
		var event helix.EventSubChannelSubscribeEvent
		if err := json.Unmarshal(notification.Event, &event); err != nil {
			log.Printf("Error unmarshaling subscribe event: %v", err)
			return
		}
		es.eventHandler("subscribe", event)
		// Add more cases as needed
	}

	w.WriteHeader(http.StatusOK)
}
