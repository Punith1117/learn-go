package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestGetEvents(t *testing.T) {
	store := NewEventStore()

	handler := getEventsHandler(store)

	req := httptest.NewRequest(
		http.MethodGet,
		"/events",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	contentType := recorder.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Fatalf(
			"expected application/json, got %q",
			contentType,
		)
	}

	var events []Event

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&events)
	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(events) != 2 {
		t.Fatalf(
			"expected 2 events, got %d",
			len(events),
		)
	}
}

func TestGetEvent(t *testing.T) {
	store := NewEventStore()

	handler := getEventHandler(store)

	req := httptest.NewRequest(
		http.MethodGet,
		"/events/event-1",
		nil,
	)

	req.SetPathValue(
		"id",
		"event-1",
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var event Event

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&event)
	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if event.ID != "event-1" {
		t.Fatalf(
			"expected event-1, got %s",
			event.ID,
		)
	}
}

func TestGetEventNotFound(t *testing.T) {
	store := NewEventStore()

	handler := getEventHandler(store)

	req := httptest.NewRequest(
		http.MethodGet,
		"/events/does-not-exist",
		nil,
	)

	req.SetPathValue(
		"id",
		"does-not-exist",
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			recorder.Code,
		)
	}
}

func TestCreateEvent(t *testing.T) {
	store := NewEventStore()

	handler := createEventHandler(store)

	body := strings.NewReader(`{
		"name": "Go First Conference",
		"available": 100
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		body,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			recorder.Code,
		)
	}

	var event Event

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&event)
	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if event.Name != "Go First Conference" {
		t.Fatalf(
			"expected Go First Conference, got %s",
			event.Name,
		)
	}

	if event.Available != 100 {
		t.Fatalf(
			"expected 100 tickets, got %d",
			event.Available,
		)
	}

	// Verify POST actually persisted the event.
	stored, exists := store.GetByID(event.ID)

	if !exists {
		t.Fatal("created event was not stored")
	}

	if stored.Name != "Go First Conference" {
		t.Fatalf(
			"expected stored event name to match, got %s",
			stored.Name,
		)
	}
}

func TestCreateEventInvalidJSON(t *testing.T) {
	store := NewEventStore()

	handler := createEventHandler(store)

	body := strings.NewReader(`invalid-json`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		body,
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestCreateEventValidation(t *testing.T) {
	store := NewEventStore()

	handler := createEventHandler(store)

	body := strings.NewReader(`{
		"name": "",
		"available": 0
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		body,
	)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestConcurrentEventCreation(t *testing.T) {
	store := NewEventStore()

	const total = 100

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()

			store.Create(
				"Concurrent Event",
				100,
			)
		}()
	}

	wg.Wait()

	events := store.GetAll()

	if len(events) != 2+total {
		t.Fatalf(
			"expected %d events, got %d",
			2+total,
			len(events),
		)
	}
}
