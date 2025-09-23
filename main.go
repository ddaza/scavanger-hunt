package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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

	// Minimal WhatsApp webhook that returns TwiML "hello world".
	// Twilio will POST x-www-form-urlencoded with fields like Body, From, To.
	mux.HandleFunc("/webhook/whatsapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("method not allowed"))
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("parse form error: %v", err)
		}

		from := r.PostFormValue("From")
		body := r.PostFormValue("Body")
		log.Printf("incoming WhatsApp message from %s: %q", from, body)

		// Respond with TwiML. WhatsApp supports TwiML responses from webhooks.
		w.Header().Set("Content-Type", "application/xml")
		// Keep it simple for now; no clue logic yet.
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>\n<Response><Message>Hello from Scavenger Hunt! 👋</Message></Response>`)
	})

	addr := ":" + port
	log.Printf("listening on %s", addr)
	// TODO: add Twilio signature validation with TWILIO_AUTH_TOKEN when moving beyond hello world.
	log.Fatal(http.ListenAndServe(addr, mux))
}
