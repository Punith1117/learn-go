package main

import (
	"fmt"
	"time"
)

func processBooking(id int, results chan string) {
	fmt.Println("Booking", id, "started")

	// Simulate different booking durations.
	time.Sleep(time.Duration(id) * time.Second)

	results <- fmt.Sprintf("Booking %d completed", id)
}

func main() {
	// Channel used to send booking results
	// from goroutines back to main.
	results := make(chan string)

	// Start all three bookings concurrently.
	go processBooking(1, results)
	go processBooking(2, results)
	go processBooking(3, results)

	// Wait for and receive three results.
	fmt.Println(<-results)
	fmt.Println(<-results)
	fmt.Println(<-results)

	fmt.Println("All bookings processed")
}
