package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type emailConfig struct {
	senderEmail string
	password    string
	smtpServer  string
	smtpPort    string
}

type emailRequest struct {
	Subject    string   `json:"subject"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
}

var client *mongo.Client

func connectToMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
}

func getAllEmailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
}

func isValidEmail(email string) bool {
	const emailRegexPattern = `^([A-Z0-9_+-]+\.?)*[A-Z0-9_+-]@([A-Z0-9][A-Z0-9-]*\.)+[A-Z]{2,}$/i`

	matched, err := regexp.MatchString(emailRegexPattern, email)
	if err != nil {
		return false
	}

	return matched
}

func getEmailConfig() (emailConfig, error) {
	config := emailConfig{
		senderEmail: os.Getenv("SENDER_EMAIL"),
		password:    os.Getenv("EMAIL_PASSWORD"),
		smtpServer:  os.Getenv("SMTP_SERVER"),
		smtpPort:    os.Getenv("SMTP_PORT"),
	}

	if config.senderEmail == "" || config.password == "" || config.smtpServer == "" || config.smtpPort == "" {
		return emailConfig{}, fmt.Errorf("one or more environment variables are not set")
	}

	if !isValidEmail(config.senderEmail) {
		return emailConfig{}, fmt.Errorf("sender email address is not valid")
	}

	return config, nil
}

func main() {
	connectToMongoDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %s", err)
		}
	}()

	http.HandleFunc("/send-email", sendEmailHandler)
	http.HandleFunc("/get-all-emails", getAllEmailsHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server start error: %s", err)
	}
}
