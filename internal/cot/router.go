package cot

import "sync/atomic"

// Router handles CoT fanout to matching subscriptions.
type Router struct {
	subs      *Subscriptions
	delivered uint64
	dropped   uint64
}

// NewRouter creates a new CoT router.
func NewRouter(subs *Subscriptions) *Router {
	if subs == nil {
		subs = NewSubscriptions()
	}
	return &Router{subs: subs}
}

// Route forwards an event to all matching subscribers.
func (r *Router) Route(evt Event) int {
	matched := r.subs.matching(evt)
	delivered := 0

	for _, sub := range matched {
		select {
		case sub.Sink <- evt:
			delivered++
			atomic.AddUint64(&r.delivered, 1)
		default:
			atomic.AddUint64(&r.dropped, 1)
		}
	}

	return delivered
}

// DeliveredCount returns total delivered messages.
func (r *Router) DeliveredCount() uint64 {
	return atomic.LoadUint64(&r.delivered)
}

// DroppedCount returns total dropped messages.
func (r *Router) DroppedCount() uint64 {
	return atomic.LoadUint64(&r.dropped)
}
