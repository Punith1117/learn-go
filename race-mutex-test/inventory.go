package main

import "sync"

type TicketInventory struct {
	mu      sync.Mutex
	tickets int
}

func NewTicketInventory(tickets int) *TicketInventory {
	return &TicketInventory{
		tickets: tickets,
	}
}

func (i *TicketInventory) Book() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.tickets <= 0 {
		return false
	}

	i.tickets--

	return true
}

func (i *TicketInventory) Remaining() int {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.tickets
}
