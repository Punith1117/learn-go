package main

import "fmt"

type Event struct {
	ID           string
	Name         string
	TotalTickets int
	Price        int
}

func (e Event) IsSoldOut() bool {
	return e.TotalTickets == 0
}

func (e *Event) Reserve() bool {
	if e.IsSoldOut() {
		return false
	}

	e.TotalTickets--
	return true
}

func main() {
	event := Event{
		ID:           "E001",
		Name:         "Coldplay",
		TotalTickets: 2,
		Price:        5000,
	}

	fmt.Println("Event:", event.Name)
	fmt.Println("Sold out:", event.IsSoldOut())

	fmt.Println("Reservation 1:", event.Reserve())
	fmt.Println("Remaining tickets:", event.TotalTickets)

	fmt.Println("Reservation 2:", event.Reserve())
	fmt.Println("Remaining tickets:", event.TotalTickets)

	fmt.Println("Reservation 3:", event.Reserve())
	fmt.Println("Remaining tickets:", event.TotalTickets)

	fmt.Println("Sold out:", event.IsSoldOut())
}
