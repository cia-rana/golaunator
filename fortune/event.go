package fortune

type EventType int

const (
	Site EventType = iota
	Circle
)

type Event struct {
	Type     EventType
	Point    *Vector
	Parabola *Parabola
}

type EventQueue struct {
	data []*Event
	Max  *Event
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		data: []*Event{&Event{}},
	}
}

func (eq *EventQueue) Push(e *Event) {
	if e == nil || e.Point == nil {
		return
	}

	if eq.Max == nil || e.Point.Y > eq.Max.Point.Y {
		eq.Max = e
	}

	eq.data = append(eq.data, e)
	for i := len(eq.data) - 1; i > 1; {
		parentI := i / 2
		if eq.data[i].Point.Y < eq.data[parentI].Point.Y || eq.data[i].Point.Y == eq.data[parentI].Point.Y && eq.data[i].Point.X < eq.data[i].Point.X {
			eq.data[i], eq.data[parentI] = eq.data[parentI], eq.data[i]
			i = parentI
		} else {
			break
		}
	}
}

func (eq *EventQueue) Pull() *Event {
	if eq.Empty() {
		return nil
	}

	result := eq.data[1]

	eq.data[1] = eq.data[len(eq.data)-1]
	eq.data = eq.data[:len(eq.data)-1]

	for childI := 2; childI < len(eq.data); {
		if childI+1 < len(eq.data) && (eq.data[childI].Point.Y > eq.data[childI+1].Point.Y || eq.data[childI].Point.Y == eq.data[childI+1].Point.Y && eq.data[childI].Point.X > eq.data[childI+1].Point.X) {
			childI++
		}

		i := childI / 2
		if eq.data[i].Point.Y > eq.data[childI].Point.Y || eq.data[i].Point.Y == eq.data[childI].Point.Y && eq.data[i].Point.X > eq.data[childI].Point.X {
			eq.data[i], eq.data[childI] = eq.data[childI], eq.data[i]
			childI *= 2
		} else {
			break
		}
	}

	if eq.Empty() {
		eq.Max = nil
	}

	return result
}

func (eq EventQueue) Empty() bool {
	return len(eq.data) <= 1
}
