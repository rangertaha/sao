package cot

import (
	"sync"
)

// Filter determines if a CoT event should be routed to a subscription.
type Filter func(Event) bool

// Subscription describes a single CoT consumer.
type Subscription struct {
	ID     string
	Filter Filter
	Sink   chan<- Event
}

// Subscriptions manages in-memory event subscribers.
type Subscriptions struct {
	mu   sync.RWMutex
	subs map[string]Subscription
}

// NewSubscriptions creates a new subscription manager.
func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		subs: make(map[string]Subscription),
	}
}

// Add registers or replaces a subscription by ID.
func (s *Subscriptions) Add(sub Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[sub.ID] = sub
}

// Remove unregisters a subscription by ID.
func (s *Subscriptions) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subs, id)
}

func (s *Subscriptions) matching(evt Event) []Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Subscription, 0, len(s.subs))
	for _, sub := range s.subs {
		if sub.Filter == nil || sub.Filter(evt) {
			result = append(result, sub)
		}
	}
	return result
}
