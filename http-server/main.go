package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Event struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Available int    `json:"available"`
}

type CreateEventRequest struct {
	Name      string `json:"name"`
	Available int    `json:"available"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// EventStore is a concurrent-safe in-memory store.
type EventStore struct {
	mu     sync.RWMutex
	events map[string]Event
	nextID int
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: map[string]Event{
			"event-1": {
				ID:        "event-1",
				Name:      "Go Conference",
				Available: 100,
			},
			"event-2": {
				ID:        "event-2",
				Name:      "Music Festival",
				Available: 500,
			},
		},
		nextID: 3,
	}
}

func (s *EventStore) GetAll() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]Event, 0, len(s.events))

	for _, event := range s.events {
		events = append(events, event)
	}

	return events
}

func (s *EventStore) GetByID(id string) (Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]

	return event, exists
}

func (s *EventStore) Create(name string, available int) Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := "event-" + strconv.Itoa(s.nextID)
	s.nextID++

	event := Event{
		ID:        id,
		Name:      name,
		Available: available,
	}

	s.events[id] = event

	return event
}

func main() {
	store := NewEventStore()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /events", getEventsHandler(store))
	mux.HandleFunc("GET /events/{id}", getEventHandler(store))
	mux.HandleFunc("POST /events", createEventHandler(store))

	var handler http.Handler = mux

	handler = loggingMiddleware(handler)
	handler = requestIDMiddleware(handler)
	handler = recoverMiddleware(handler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Println("server running on http://localhost:8080")

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}

func getEventsHandler(store *EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events := store.GetAll()

		writeJSON(w, http.StatusOK, events)
	}
}

func getEventHandler(store *EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		event, exists := store.GetByID(id)

		if !exists {
			writeError(
				w,
				http.StatusNotFound,
				"event not found",
			)
			return
		}

		writeJSON(w, http.StatusOK, event)
	}
}

func createEventHandler(store *EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEventRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(
				w,
				http.StatusBadRequest,
				"invalid JSON",
			)
			return
		}

		req.Name = strings.TrimSpace(req.Name)

		if req.Name == "" {
			writeError(
				w,
				http.StatusBadRequest,
				"name is required",
			)
			return
		}

		if req.Available <= 0 {
			writeError(
				w,
				http.StatusBadRequest,
				"available must be greater than zero",
			)
			return
		}

		event := store.Create(
			req.Name,
			req.Available,
		)

		writeJSON(
			w,
			http.StatusCreated,
			event,
		)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Error: message,
		},
	)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {
			requestID = strconv.FormatInt(
				time.Now().UnixNano(),
				10,
			)
		}

		w.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(
					"panic recovered: %v",
					err,
				)

				writeError(
					w,
					http.StatusInternalServerError,
					"internal server error",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
