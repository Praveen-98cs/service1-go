package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	service2URL := os.Getenv("SERVICE2_URL")
	if service2URL == "" {
		log.Fatal("SERVICE2_URL environment variable is required")
	}

	http.HandleFunc("/service1", func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("x-request-id")
		if requestID != "" {
			fmt.Println("Received request with x-request-id:", requestID)
		}

		req, err := http.NewRequest("GET", service2URL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if requestID != "" {
			req.Header.Set("x-request-id", requestID)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	fmt.Printf("Service1 listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
