package main

import "fmt"

type Event struct {
	ID           string
	Name         string
	TotalTickets int
	Price        int
}

func main() {
	events := []Event{}

	events = append(events, Event{
		ID:           "E001",
		Name:         "Coldplay",
		TotalTickets: 50000,
		Price:        5000,
	})

	events = append(events, Event{
		ID:           "E002",
		Name:         "IPL Final",
		TotalTickets: 30000,
		Price:        2500,
	})

	events = append(events, Event{
		ID:           "E003",
		Name:         "Tech Conference",
		TotalTickets: 5000,
		Price:        1000,
	})

	for _, event := range events {
		fmt.Println(event.ID, event.Name, event.TotalTickets, event.Price)
	}

	eventsMap := make(map[string]Event)

	eventsMap["E001"] = Event{
		ID:           "E001",
		Name:         "Coldplay",
		TotalTickets: 50000,
		Price:        5000,
	}

	eventsMap["E002"] = events[1]
	eventsMap["E003"] = events[2]

	event002, ok := eventsMap["E002"]

	if ok {
		fmt.Println(event002)
	}
}
