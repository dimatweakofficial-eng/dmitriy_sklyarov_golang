package event

const EventLinkVisited = "link.visited"

// макет информации которую будет содержать канал с инфо о событиях
type Event struct {
	Type string
	Data any
}

// макет самого канала
type EventBus struct {
	bus chan Event
}

// макет инициализации канала
func NewEventBus() *EventBus {
	return &EventBus{
		make(chan Event),
	}
}

// передача каких то данных по макету в канал
func (e *EventBus) Publish(data Event) {
	e.bus <- data
}

// доставание информации в виде макета с канала
func (e *EventBus) Subscribe() <-chan Event {
	return e.bus
}
