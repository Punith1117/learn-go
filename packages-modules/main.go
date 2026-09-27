package main

import (
	"fmt"

	"github.com/punith1117/learn-go/packages-modules/event"
)

func main() {
	coldplay := event.NewEvent(
		"E001",
		"Coldplay",
		50000,
	)

	fmt.Println("ID:", coldplay.ID)
	fmt.Println("Name:", coldplay.Name)
	fmt.Println("Tickets:", coldplay.TotalTickets)
	fmt.Println("Sold out:", coldplay.IsSoldOut())
}
