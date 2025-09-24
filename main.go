package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"scavenger-hunt/internal/middleware"

	"github.com/twilio/twilio-go/twiml"
)

type Clue struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

const CLUES_CONFIG = "config/clues.json"

func LoadClues() ([]Clue, error) {
	f, err := os.Open(CLUES_CONFIG)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var clues []Clue
	if err := json.NewDecoder(f).Decode(&clues); err != nil {
		return nil, err
	}
	return clues, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	clues, err := LoadClues()
	if err != nil {
		log.Printf("note: could not load clues (config/clues.json): %v", err)
	} else {
		log.Printf("loaded %d clues (not used yet)", len(clues))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Minimal WhatsApp webhook protected by Twilio signature validation.
	mux.Handle("/webhook/whatsapp", middleware.TwilioAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("method not allowed"))
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("parse form error: %v", err)
		}

		message := &twiml.MessagingMessage{Body: "Yay, valid requests!"}
		twimlResult, err := twiml.Messages([]twiml.Element{message})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		from := r.PostFormValue("From")
		body := r.PostFormValue("Body")
		log.Printf("incoming WhatsApp message from %s: %q", from, body)

		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(twimlResult))
	})))

	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
