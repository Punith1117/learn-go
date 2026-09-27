package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 1. DOMAIN
// ============================================================

type Event struct {
	ID           string
	Name         string
	TotalTickets int
}

// ============================================================
// 2. REPOSITORY INTERFACE
// ============================================================

// EventRepository defines what the application needs
// from an event repository.
//
// The service doesn't care how the data is stored.
//
// Any type that has this method:
//
//	FindByID(id string) (Event, error)
//
// automatically satisfies EventRepository.
type EventRepository interface {
	FindByID(id string) (Event, error)
}

// ============================================================
// 3. IN-MEMORY REPOSITORY
// ============================================================

// InMemoryEventRepository represents a repository where
// events are stored directly in application memory.
//
// In a real application, this could be useful for:
//   - tests
//   - prototypes
//   - local development
type InMemoryEventRepository struct {
	events map[string]Event
}

// FindByID looks for an event in the in-memory map.
func (r InMemoryEventRepository) FindByID(id string) (Event, error) {
	event, exists := r.events[id]

	if !exists {
		return Event{}, errors.New("event not found")
	}

	return event, nil
}

// ============================================================
// 4. SIMULATED POSTGRESQL REPOSITORY
// ============================================================

// PostgreSQLEventRepository represents a repository that
// would normally communicate with PostgreSQL.
//
// For now, we simulate PostgreSQL using a map.
//
// Later, when we learn database access, this map will be
// replaced by an actual *sql.DB connection.
type PostgreSQLEventRepository struct {
	// Pretend this data lives inside PostgreSQL.
	events map[string]Event
}

// FindByID simulates a PostgreSQL query.
//
// A real implementation would eventually do something like:
//
//	SELECT id, name, total_tickets
//	FROM events
//	WHERE id = $1
//
// For now, we perform the same lookup against our simulated
// database map.
func (r PostgreSQLEventRepository) FindByID(id string) (Event, error) {
	fmt.Println("[PostgreSQL] Executing query for event:", id)

	event, exists := r.events[id]

	if !exists {
		return Event{}, errors.New("event not found in PostgreSQL")
	}

	return event, nil
}

// ============================================================
// 5. SERVICE
// ============================================================

// EventService contains application/business logic.
//
// Notice that it depends on:
//
//	EventRepository
//
// It does NOT depend directly on:
//
//	InMemoryEventRepository
//	PostgreSQLEventRepository
//
// This is the important part of the design.
type EventService struct {
	repository EventRepository
}

// NewEventService creates an EventService.
//
// The parameter is an interface:
//
//	EventRepository
//
// Therefore we can pass either:
//
//	InMemoryEventRepository
//
// or:
//
//	PostgreSQLEventRepository
func NewEventService(repository EventRepository) EventService {
	return EventService{
		repository: repository,
	}
}

// GetEvent asks the repository to find an event.
//
// The service doesn't know whether the repository uses:
//
//	map
//	PostgreSQL
//	Redis
//	HTTP API
//	etc.
//
// It only knows that FindByID exists.
func (s EventService) GetEvent(id string) (Event, error) {
	return s.repository.FindByID(id)
}

// ============================================================
// 6. HELPER FUNCTION
// ============================================================

// printEvent is just a small helper so that main()
// stays easier to read.
func printEvent(event Event) {
	fmt.Println("ID:", event.ID)
	fmt.Println("Name:", event.Name)
	fmt.Println("Tickets:", event.TotalTickets)
}

// ============================================================
// 7. MAIN
// ============================================================

func main() {
	// ========================================================
	// PART 1: IN-MEMORY REPOSITORY
	// ========================================================

	fmt.Println("========== IN-MEMORY REPOSITORY ==========")

	// This map represents our in-memory data.
	memoryEvents := map[string]Event{
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

	// Create the in-memory repository.
	memoryRepository := InMemoryEventRepository{
		events: memoryEvents,
	}

	// Give the repository to the service.
	//
	// memoryRepository satisfies EventRepository,
	// so Go allows us to pass it here.
	memoryService := NewEventService(memoryRepository)

	// --------------------------------------------------------
	// Find E001
	// --------------------------------------------------------

	event, err := memoryService.GetEvent("E001")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:")
		printEvent(event)
	}

	// --------------------------------------------------------
	// Find an event that doesn't exist
	// --------------------------------------------------------

	event, err = memoryService.GetEvent("E999")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:")
		printEvent(event)
	}

	// ========================================================
	// PART 2: POSTGRESQL REPOSITORY
	// ========================================================

	fmt.Println()
	fmt.Println("========== POSTGRESQL REPOSITORY ==========")

	// Imagine these records are actually stored in PostgreSQL.
	//
	// We're using a map only to simulate the database while
	// we're learning the Go concepts.
	postgresEvents := map[string]Event{
		"E100": {
			ID:           "E100",
			Name:         "Taylor Swift",
			TotalTickets: 60000,
		},
		"E200": {
			ID:           "E200",
			Name:         "Rock Festival",
			TotalTickets: 40000,
		},
	}

	// Create the simulated PostgreSQL repository.
	postgresRepository := PostgreSQLEventRepository{
		events: postgresEvents,
	}

	// Create another service.
	//
	// Notice:
	//
	// The EventService itself hasn't changed.
	//
	// We're simply giving it a different repository.
	postgresService := NewEventService(postgresRepository)

	// --------------------------------------------------------
	// Find E100
	// --------------------------------------------------------

	event, err = postgresService.GetEvent("E100")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:")
		printEvent(event)
	}

	// --------------------------------------------------------
	// Find an event that doesn't exist
	// --------------------------------------------------------

	event, err = postgresService.GetEvent("E999")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Event found:")
		printEvent(event)
	}

	// ========================================================
	// PART 3: SAME SERVICE, DIFFERENT REPOSITORIES
	// ========================================================

	fmt.Println()
	fmt.Println("========== SAME SERVICE, DIFFERENT REPOSITORIES ==========")

	fmt.Println("Using in-memory repository:")
	event, err = memoryService.GetEvent("E002")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		printEvent(event)
	}

	fmt.Println()

	fmt.Println("Using PostgreSQL repository:")
	event, err = postgresService.GetEvent("E200")

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		printEvent(event)
	}
}
