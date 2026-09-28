package main

import (
	"sync"
	"testing"
)

func TestConcurrentBooking(t *testing.T) {
	inventory := NewTicketInventory(100)

	var wg sync.WaitGroup

	successfulBookings := 0

	var mu sync.Mutex

	wg.Add(1000)

	for n := 0; n < 1000; n++ {
		go func() {
			defer wg.Done()

			if inventory.Book() {
				mu.Lock()
				successfulBookings++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if successfulBookings != 100 {
		t.Fatalf(
			"expected 100 successful bookings, got %d",
			successfulBookings,
		)
	}

	if inventory.Remaining() != 0 {
		t.Fatalf(
			"expected 0 tickets remaining, got %d",
			inventory.Remaining(),
		)
	}
}
