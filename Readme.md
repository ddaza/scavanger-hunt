# Scavenger Hunt Service (Go + Twilio WhatsApp)

A minimal Go webhook server for Twilio WhatsApp that currently returns a TwiML "hello world". It also includes a placeholder JSON data source for future clue logic.

## Quick Start

- Requirements: Go 1.22+
- Local run:
  - `cp .env.example .env` (optional, to change `PORT`)
  - `make run` or `go run .`
  - Health check: `curl -s http://localhost:8080/healthz`

## WhatsApp Webhook (Hello World)

- Endpoint: `POST /webhook/whatsapp`
- Expected content type: `application/x-www-form-urlencoded`
- Returns: TwiML with a simple message

Example local test:

```
curl -X POST http://localhost:8080/webhook/whatsapp \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data 'From=whatsapp:+14155238886&To=whatsapp:+1234567890&Body=hello'
```

You should receive a TwiML response like:

```
<?xml version="1.0" encoding="UTF-8"?>
<Response><Message>Hello from Scavenger Hunt! 👋</Message></Response>
```

## Twilio WhatsApp Setup (dev sandbox)

1. Enable the WhatsApp Sandbox in Twilio.
2. Set the sandbox "When a message comes in" webhook to your public URL (e.g., via `ngrok` or Cloudflare Tunnel):
   - `https://<your-public-host>/webhook/whatsapp`
3. Send a WhatsApp message to your Twilio sandbox number; Twilio will POST to the webhook and relay the TwiML reply back to the user.

> Note: For production, add Twilio signature validation using `TWILIO_AUTH_TOKEN` and HTTPS.

## Project Layout

- `main.go` – HTTP server and WhatsApp webhook handler.
- `config/clues.json` – placeholder data source for future clue logic.
- `.env.example` – environment variables.
- `Makefile` – convenience `run` and `build` targets.

## Deploy on Render (Web Service)

- Commit and push this repo to GitHub/GitLab.
- This repo includes `render.yaml:1` so you can either:
  - Use Render Blueprints: `New` → `Blueprint` → select the repo (it will create a Web Service per `render.yaml`).
  - Or create manually as a Web Service:
    - Environment: `Go`
    - Build Command: `go build -o bin/scavenger-hunt .`
    - Start Command: `./bin/scavenger-hunt`
    - Health Check Path: `/healthz`
    - Auto deploy on push: enabled

Environment variables (Render Dashboard → Environment):
- `TWILIO_AUTH_TOKEN` (optional now; required when you enable signature validation)
- `PORT` is provided by Render automatically; the server already respects it.

After deploy, set your Twilio WhatsApp webhook to:
- `https://<your-service>.onrender.com/webhook/whatsapp`

## Next Steps (not implemented yet)

- Load clues from JSON and track per-user progress (keyed by `From`).
- Validate Twilio signatures on incoming requests.
- Add persistence (DB) once JSON prototype is validated.
