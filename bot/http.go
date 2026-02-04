package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// CommandHandler handles /command endpoint.
func (b *Bot) CommandHandler(verify bool, trustedIssuers map[string]struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		botID, ok := getUUID(raw, "bot_id")
		if !ok {
			http.Error(w, "missing bot_id", http.StatusBadRequest)
			return
		}
		if verify {
			if err := b.VerifyRequest(botID, r.Header.Get("Authorization"), trustedIssuers); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		command, _ := raw["command"].(map[string]any)
		cmdType, _ := command["command_type"].(string)

		if cmdType == "system" || strings.HasPrefix(commandString(command), "system:") {
			ev, err := ParseSystemEvent(raw)
			if err == nil {
				go b.HandleSystemEvent(context.Background(), ev)
			}
		} else {
			msg, err := ParseIncomingMessage(raw)
			if err == nil {
				go b.HandleIncomingMessage(context.Background(), msg)
			}
		}

		writeJSON(w, BuildCommandAcceptedResponse(), http.StatusAccepted)
	}
}

// StatusHandler handles /status endpoint.
func (b *Bot) StatusHandler(verify bool, trustedIssuers map[string]struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := map[string]string{}
		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				q[k] = v[0]
			}
		}

		recipient, err := ParseStatusRecipient(q)
		if err != nil {
			http.Error(w, "invalid status params", http.StatusBadRequest)
			return
		}
		if verify {
			if err := b.VerifyRequest(recipient.BotID, r.Header.Get("Authorization"), trustedIssuers); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		menu, err := b.Status(r.Context(), recipient)
		if err != nil {
			http.Error(w, "status error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]any{"status": "ok", "commands": menu}, http.StatusOK)
	}
}

// CallbackHandler handles /notification/callback endpoint.
func (b *Bot) CallbackHandler(verify bool, trustedIssuers map[string]struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		botID, ok := getUUID(raw, "bot_id")
		if !ok {
			http.Error(w, "missing bot_id", http.StatusBadRequest)
			return
		}
		if verify {
			if err := b.VerifyRequest(botID, r.Header.Get("Authorization"), trustedIssuers); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		cb, err := ParseCallback(raw)
		if err == nil {
			_ = b.SetCallback(cb)
		}
		writeJSON(w, BuildCommandAcceptedResponse(), http.StatusAccepted)
	}
}

// SyncSmartAppHandler handles /smartapps/request endpoint.
func (b *Bot) SyncSmartAppHandler(verify bool, trustedIssuers map[string]struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		botID, ok := getUUID(raw, "bot_id")
		if !ok {
			http.Error(w, "missing bot_id", http.StatusBadRequest)
			return
		}
		if verify {
			if err := b.VerifyRequest(botID, r.Header.Get("Authorization"), trustedIssuers); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		ev, err := ParseSmartAppEvent(raw)
		if err != nil {
			http.Error(w, "invalid event", http.StatusBadRequest)
			return
		}
		resp, err := b.HandleSyncSmartApp(r.Context(), ev)
		if err != nil {
			writeJSON(w, map[string]any{"status": "error", "reason": "smartapp_error"}, http.StatusOK)
			return
		}
		writeJSON(w, resp, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func commandString(cmd map[string]any) string {
	if cmd == nil {
		return ""
	}
	if body, ok := cmd["body"].(string); ok {
		return body
	}
	return ""
}
