package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type GuardMetric struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Level     int       `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	activeMetrics []GuardMetric
	mu            sync.RWMutex
	metricPool    = sync.Pool{
		New: func() interface{} {
			return new(GuardMetric)
		},
	}
)

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/ingest", handleIngest)
	mux.HandleFunc("/api/v1/metrics", handleGetMetrics)
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("System is up"))
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server started on 8080")
	log.Fatal(srv.ListenAndServe())
}

func handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	m := metricPool.Get().(*GuardMetric)
	defer metricPool.Put(m)

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	m.Timestamp = time.Now()
	
	mu.Lock()
	activeMetrics = append(activeMetrics, *m)
	if len(activeMetrics) > 5000 {
		activeMetrics = activeMetrics[1:]
	}
	mu.Unlock()
	
	w.WriteHeader(http.StatusCreated)
}

func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	data, err := json.Marshal(activeMetrics)
	mu.RUnlock()

	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
