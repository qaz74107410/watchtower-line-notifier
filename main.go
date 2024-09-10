package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type LineNotifier struct {
	APIEndpoint string
	Token       string
}

type WatchtowerPayload struct {
	Text string `json:"text"`
}

func main() {
	// Initialize configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %s", err)
	}

	// Create LineNotifier instance
	notifier := &LineNotifier{
		APIEndpoint: viper.GetString("line.api_endpoint"),
		Token:       viper.GetString("line.token"),
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, notifier)
	})

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; copy from example
			log.Println("Config file not found. Copying from example...")
			if err := copyFile("config.yaml.example", "config.yaml"); err != nil {
				return fmt.Errorf("error copying config file: %s", err)
			}
			// Now try to read the config again
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading copied config file: %s", err)
			}
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %s", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request, notifier *LineNotifier) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the raw body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Convert body to string
	messageContent := string(body)

	// You might want to limit the message length to avoid very large messages
	if len(messageContent) > 1000 {
		messageContent = messageContent[:1000] + "... (message truncated)"
	}

	message := fmt.Sprintf("Watchtower Update: %s", messageContent)
	fmt.Println(message)
	err = notifier.SendLineNotification(message)
	if err != nil {
		http.Error(w, "Error sending LINE notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent to LINE"))
}

func (ln *LineNotifier) SendLineNotification(message string) error {
	form := url.Values{}
	form.Add("message", message)

	req, err := http.NewRequest("POST", ln.APIEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+ln.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}
