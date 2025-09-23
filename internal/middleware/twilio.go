package middleware

import (
    "bytes"
    "io"
    "mime"
    "net/http"
    "os"
    "sort"
    "strings"

    twclient "github.com/twilio/twilio-go/client"
)

// TwilioAuth returns middleware that validates Twilio's X-Twilio-Signature
// header using the configured Auth Token. It supports both form-encoded and
// JSON payloads (Validate vs ValidateBody). If validation fails, it returns 403.
func TwilioAuth(next http.Handler) http.Handler {
    authToken := os.Getenv("TWILIO_AUTH_TOKEN")
    publicBase := strings.TrimRight(os.Getenv("PUBLIC_BASE_URL"), "/")
    validator := twclient.NewRequestValidator(authToken)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if authToken == "" {
            // If no token is set, skip validation (useful for local dev).
            next.ServeHTTP(w, r)
            return
        }

        sig := r.Header.Get("X-Twilio-Signature")
        if sig == "" {
            http.Error(w, "missing signature", http.StatusForbidden)
            return
        }

        // Build the full URL Twilio used, accounting for proxies (Render) and optional PUBLIC_BASE_URL.
        url := buildFullURL(r, publicBase)

        ctype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
        switch ctype {
        case "application/x-www-form-urlencoded":
            // Parse and build a flat map[string]string (first value per key).
            if err := r.ParseForm(); err != nil {
                http.Error(w, "bad form", http.StatusBadRequest)
                return
            }
            params := make(map[string]string, len(r.PostForm))
            // Twilio requires sorted parameters; the validator handles sorting internally,
            // but we normalize consistently here.
            keys := make([]string, 0, len(r.PostForm))
            for k := range r.PostForm {
                keys = append(keys, k)
            }
            sort.Strings(keys)
            for _, k := range keys {
                if vs := r.PostForm[k]; len(vs) > 0 {
                    params[k] = vs[0]
                }
            }
            if !validator.Validate(url, params, sig) {
                http.Error(w, "invalid signature", http.StatusForbidden)
                return
            }
        default:
            // Read raw body for JSON or other content types, then restore it for downstream handlers.
            body, err := io.ReadAll(r.Body)
            if err != nil {
                http.Error(w, "bad body", http.StatusBadRequest)
                return
            }
            r.Body.Close()
            r.Body = io.NopCloser(bytes.NewReader(body))

            if !validator.ValidateBody(url, body, sig) {
                http.Error(w, "invalid signature", http.StatusForbidden)
                return
            }
        }

        next.ServeHTTP(w, r)
    })
}

func buildFullURL(r *http.Request, publicBase string) string {
    // Prefer explicit base when provided to avoid proxy/scheme mismatches.
    base := publicBase
    if base == "" {
        scheme := r.Header.Get("X-Forwarded-Proto")
        if scheme == "" {
            if r.TLS != nil {
                scheme = "https"
            } else {
                scheme = "http"
            }
        }
        host := r.Header.Get("X-Forwarded-Host")
        if host == "" {
            host = r.Host
        }
        base = scheme + "://" + host
    }
    if rq := r.URL.RawQuery; rq != "" {
        return base + r.URL.Path + "?" + rq
    }
    return base + r.URL.Path
}
