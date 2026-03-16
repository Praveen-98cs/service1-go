package main

import (
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Starting service1...")
	log.Printf("SERVICE2_URL: %s", service2URL)
	log.Printf("PORT: %s", port)

	http.HandleFunc("/service1", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		log.Printf("Incoming headers: %v", r.Header)

		requestID := r.Header.Get("x-request-id")
		if requestID != "" {
			log.Printf("x-request-id: %s", requestID)
		} else {
			log.Printf("WARNING: No x-request-id header in request")
		}

		log.Printf("Forwarding GET request to: %s", service2URL)

		req, err := http.NewRequest("GET", service2URL, nil)
		if err != nil {
			log.Printf("ERROR creating request to service2: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if requestID != "" {
			req.Header.Set("x-request-id", requestID)
		}

		log.Printf("Outgoing headers to service2: %v", req.Header)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("ERROR calling service2: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		log.Printf("Response from service2: status=%d", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR reading service2 response body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Response body from service2: %s", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
		log.Printf("Request completed successfully")
	})

	log.Printf("Service1 listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
