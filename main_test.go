package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestWebhookHandler(t *testing.T) {
	// Set up a mock LINE Notify server
	mockLineServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("Expected Authorization header with token, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockLineServer.Close()

	// Set up Viper configuration
	viper.Set("line.api_endpoint", mockLineServer.URL)
	viper.Set("line.token", "testtoken")

	// Create LineNotifier instance
	notifier := &LineNotifier{
		APIEndpoint: viper.GetString("line.api_endpoint"),
		Token:       viper.GetString("line.token"),
	}

	// Create a test server with our webhookHandler
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, notifier)
	}))
	defer ts.Close()

	// Test case
	payload := WatchtowerPayload{Text: "Test update"}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(ts.URL+"/webhook", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestSendLineNotification(t *testing.T) {
	// Set up a mock LINE Notify server
	mockLineServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("Expected Authorization header with token, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockLineServer.Close()

	notifier := &LineNotifier{
		APIEndpoint: mockLineServer.URL,
		Token:       "testtoken",
	}

	err := notifier.SendLineNotification("Test message")
	if err != nil {
		t.Errorf("SendLineNotification failed: %v", err)
	}
}
