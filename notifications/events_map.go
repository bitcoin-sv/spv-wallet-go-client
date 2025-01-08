package notifications

import "sync"

type eventsMap struct {
	registered *sync.Map
}

func newEventsMap() *eventsMap {
	return &eventsMap{
		registered: &sync.Map{},
	}
}

func (em *eventsMap) store(name string, handler *eventHandler) {
	em.registered.Store(name, handler)
}

func (em *eventsMap) load(name string) (*eventHandler, bool) {
	h, ok := em.registered.Load(name)
	if !ok {
		return nil, false
	}
	return h.(*eventHandler), true
}
