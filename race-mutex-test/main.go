package main

import (
	"fmt"
	"sync"
)

func main() {
	inventory := NewTicketInventory(100)

	var wg sync.WaitGroup

	wg.Add(1000)

	for n := 0; n < 1000; n++ {
		go func() {
			defer wg.Done()

			inventory.Book()
		}()
	}

	wg.Wait()

	fmt.Println("Tickets remaining:", inventory.Remaining())
}
