package main

import (
	"context"
	"fmt"
	"time"
)

func processBooking(ctx context.Context) error {
	fmt.Println("Booking started")

	// Always runs when processBooking returns,
	// whether it succeeds or returns an error.
	defer fmt.Println("Booking cleanup")

	// Simulate a booking operation that takes 3 seconds.
	select { // Run whichever channel that completes first
	case <-time.After(3 * time.Second):
		fmt.Println("Booking completed")
		return nil

	case <-ctx.Done():
		fmt.Println("Booking cancelled")
		return ctx.Err()
	}
}

func main() {
	// The booking is allowed to run for at most 2 seconds.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	err := processBooking(ctx)
	if err != nil {
		fmt.Println("Booking failed:", err)
		return
	}

	fmt.Println("Booking successful")
}
