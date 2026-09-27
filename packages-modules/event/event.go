package event

type Event struct {
	ID           string
	Name         string
	TotalTickets int
}

func NewEvent(id string, name string, totalTickets int) Event {
	return Event{
		ID:           id,
		Name:         name,
		TotalTickets: totalTickets,
	}
}

func (e Event) IsSoldOut() bool {
	return e.TotalTickets == 0
}
