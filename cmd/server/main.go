package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bookshelf/monolith/internal/config"
)

func main() {
	config := config.Load()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := map[string]string{
			"status": "ok",
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Println("failed to encode response:", err)
		}
	})

	if err := http.ListenAndServe(":" + config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
