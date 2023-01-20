package queue

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

type StandardQueue struct {
	Id       string
	Elements StandardCustomers
}

func NewStandardQueue(ctx context.Context, id string) (*StandardQueue, error) {

	return &StandardQueue{
		Id:       id,
		Elements: make(StandardCustomers, 0),
	}, nil
}

func (q *StandardQueue) Add(c *customer.StandardCustomer) {
	currQueue := q.Elements
	currQueue = append(currQueue, c)

	q.Elements = currQueue
}

func (q *StandardQueue) Pop() (*customer.StandardCustomer, error) {
	currQueue := q.Elements
	len := len(currQueue)
	if len == 0 {
		return nil, fmt.Errorf("cannot pop from empty queue")
	}

	customer := currQueue[len-1]
	q.Elements = currQueue[:len]

	return customer, nil
}

func (q *StandardQueue) Len() int {
	return len(q.Elements)
}
