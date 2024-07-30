package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

type emailConfig struct {
	senderEmail string
	password    string
	smtpServer  string
	smtpPort    string
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
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
	http.HandleFunc("/send-email", sendEmailHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server start error: %s", err)
	}
}
