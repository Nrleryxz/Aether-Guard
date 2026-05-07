package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	AIScore   int       `json:"ai_score"`
	AINote    string    `json:"ai_note"`
	Timestamp time.Time `json:"timestamp"`
}

type AIRequest struct {
	Message string `json:"message"`
	Level   int    `json:"level"`
}

type AIResponse struct {
	Score int    `json:"score"`
	Note  string `json:"note"`
}

var (
	activeMetrics []GuardMetric
	mu            sync.RWMutex
	metricPool    = sync.Pool{
		New: func() interface{} {
			return new(GuardMetric)
		},
	}
	client = &http.Client{
		Timeout: 2 * time.Second,
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

	aiRes, err := analyzeWithAI(m.Message, m.Level)
	if err == nil {
		m.AIScore = aiRes.Score
		m.AINote = aiRes.Note
	} else {
		m.AINote = "AI Analysis Failed"
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

func analyzeWithAI(message string, level int) (*AIResponse, error) {
	apiURL := "http://localhost:5000/analyze"
	reqBody, _ := json.Marshal(AIRequest{Message: message, Level: level})
	
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var aiRes AIResponse
	if err := json.Unmarshal(body, &aiRes); err != nil {
		return nil, err
	}

	return &aiRes, nil
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}
