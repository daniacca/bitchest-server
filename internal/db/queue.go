package db

import "errors"

type Queue struct {
	items []string
	length int
}

func NewQueue() *Queue {
	return &Queue{
		items: []string{},
		length: 0,
	}
}

// Push adds an item to the end of the queue
func (q *Queue) Push(item string) {
	q.items = append(q.items, item)
	q.length++
}

// Pop removes and returns the last item of the queue
func (q *Queue) Pop() (string, error) {
	if q.length == 0 {
		return "", errors.New("queue is empty")
	}
	item := q.items[q.length-1]
	q.items = q.items[:q.length-1]
	q.length--
	return item, nil
}

// Shift removes and returns the first item of the queue
func (q *Queue) Shift() (string, error) {
	if q.length == 0 {
		return "", errors.New("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	q.length--
	return item, nil
}

// Unshift adds an item to the beginning of the queue
func (q *Queue) Unshift(item string) {
	q.items = append([]string{item}, q.items...)
	q.length++
}

// GetLength returns the number of items in the queue
func (q Queue) GetLength() int {
	return q.length
}

// GetItems returns the items in the queue
func (q Queue) GetItems() []string {
	return q.items
}

// GetSize returns the size (in bytes) of the queue
func (q Queue) GetSize() int {
	sum := 0
	for _, s := range q.items {
		sum += len(s)
	}
	return sum
}

func (q Queue) Index(idx int) (string, error) {
    n := q.GetLength()
    if idx < 0 || idx >= n {
        return "", errors.New("index out of range")
    }

    return q.items[idx], nil
}

// Set sets the item at the given index 
func (q *Queue) Set(idx int, value string) error {
    n := q.GetLength()
    if idx < 0 || idx >= n {
        return errors.New("index out of range")
    }

    q.items[idx] = value
    return nil
}