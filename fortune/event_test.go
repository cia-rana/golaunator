package fortune

import (
	"testing"
)

func TestEventQueueEmpty(t *testing.T) {
	eq := NewEventQueue()
	if !eq.Empty() {
		t.Fatal("Event-Queue should be empty.")
	}

	eq.Push(&Event{Point: &Vector{}})
	if eq.Empty() {
		t.Fatal("Event-Queue should not be empty.")
	}
}

func TestEventQueue(t *testing.T) {
	eq := NewEventQueue()

	data := []float64{1, 9, 3, 7, 5, 6, 4, 8, 2, 10}
	expected := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, d := range data {
		eq.Push(&Event{Point: &Vector{X: 0, Y: d}})
	}

	for i, d := range expected {
		e := eq.Pull()
		if e == nil {
			t.Fatal()
		} else if d != e.Point.Y {
			t.Fatalf("no: %d, expected: %f, actual %f", i, d, e.Point.Y)
		}
	}
}
