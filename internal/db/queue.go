package db

type Queue struct {
	items []string
	length int
}

// Push adds an item to the end of the queue
func (q *Queue) Push(item string) {
	q.items = append(q.items, item)
	q.length++
}

// Pop removes and returns the last item of the queue
func (q *Queue) Pop() string {
	item := q.items[q.length-1]
	q.items = q.items[:q.length-1]
	q.length--
	return item
}

// Shift removes and returns the first item of the queue
func (q *Queue) Shift() string {
	item := q.items[0]
	q.items = q.items[1:]
	q.length--
	return item
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