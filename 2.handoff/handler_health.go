package main

import "net/http"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, r.Context().Value(requestIDKey).(string), map[string]string{"status": "ok"})
}
