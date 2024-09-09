package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("SERVER_PORT")
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		name, ok := payload["name"].(string)
		w.WriteHeader(http.StatusOK)
		if ok {
			json.NewEncoder(w).Encode(map[string]any{"message": fmt.Sprintf("Hello, %s", name)})
		} else {
			json.NewEncoder(w).Encode(map[string]any{"message": "Hello world"})
		}
	})
	http.ListenAndServe(":"+addr, nil)
}

