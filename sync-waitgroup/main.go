package main

import (
	"fmt"
	"sync"
	"time"
)

func processBooking(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Booking", id, "started")

	// Simulate booking work.
	time.Sleep(time.Duration(id) * time.Second)

	fmt.Println("Booking", id, "completed")
}

func main() {
	var wg sync.WaitGroup

	// We are starting 3 goroutines.
	wg.Add(3)

	go processBooking(1, &wg)
	go processBooking(2, &wg)
	go processBooking(3, &wg)

	// Wait until all 3 goroutines call Done().
	wg.Wait()

	fmt.Println("All bookings processed")
}
