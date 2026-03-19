package cot

import "testing"

func TestRouterRoute(t *testing.T) {
	t.Parallel()

	subs := NewSubscriptions()
	ch := make(chan Event, 1)
	subs.Add(Subscription{
		ID:   "all",
		Sink: ch,
	})

	router := NewRouter(subs)
	delivered := router.Route(Event{UID: "A-1", Type: "a-f-G-U-C"})
	if delivered != 1 {
		t.Fatalf("expected delivered=1, got %d", delivered)
	}
}
