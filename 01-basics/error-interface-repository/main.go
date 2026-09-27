package main

import (
	"errors"
	"fmt"
)

type Event struct {
	ID           string
	Name         string
	TotalTickets int
}

// EventRepository defines what an event repository must be able to do.
type EventRepository interface {
	FindByID(id string) (Event, error)
}

// InMemoryEventRepository is an in-memory implementation
// of EventRepository.
type InMemoryEventRepository struct {
	events map[string]Event
}

// FindByID looks up an event by ID.
//
// It returns:
//   - the Event if found
//   - an error if the Event doesn't exist
func (r InMemoryEventRepository) FindByID(id string) (Event, error) {
	event, ok := r.events[id]

	if !ok {
		return Event{}, errors.New("event not found")
	}

	return event, nil
}

// EventService contains the business logic related to events.
type EventService struct {
	repository EventRepository
}

// NewEventService creates a new EventService.
func NewEventService(repository EventRepository) EventService {
	return EventService{
		repository: repository,
	}
}

// GetEvent retrieves an event through the repository.
func (s EventService) GetEvent(id string) (Event, error) {
	return s.repository.FindByID(id)
}

func main() {
	// Create some events.
	events := map[string]Event{
		"E001": {
			ID:           "E001",
			Name:         "Coldplay",
			TotalTickets: 50000,
		},
		"E002": {
			ID:           "E002",
			Name:         "IPL Final",
			TotalTickets: 30000,
		},
	}

	// Create the repository.
	repository := InMemoryEventRepository{
		events: events,
	}

	// EventService depends on the EventRepository interface,
	// not directly on InMemoryEventRepository.
	service := NewEventService(repository)

	// -------------------------
	// Successful lookup
	// -------------------------

	event, err := service.GetEvent("E001")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:")
		fmt.Println("ID:", event.ID)
		fmt.Println("Name:", event.Name)
		fmt.Println("Tickets:", event.TotalTickets)
	}

	// -------------------------
	// Failed lookup
	// -------------------------

	event, err = service.GetEvent("E999")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:", event)
	}
}
