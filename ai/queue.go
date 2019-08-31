package ai

type queue struct {
	head *State
	tail *State
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) enqueue(s *State) {
	if q.head == nil {
		q.head = s
		q.tail = s
	} else {
		q.tail.Next = s
		q.tail = s
	}
	s.Next = nil
}

func (q *queue) dequeue() *State {
	s := q.head
	if q.head != nil {
		if q.head == q.tail {
			q.head = nil
			q.tail = nil
		} else {
			q.head = q.head.Next
		}
	}
	return s
}

func (q *queue) isEmpty() bool {
	return q.head == nil
}

func (q *queue) isNotEmpty() bool {
	return q.head != nil
}
