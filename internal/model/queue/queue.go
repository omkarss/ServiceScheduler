package queue

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

type QueueType string
type CustomerQueue []*customer.Customer

const (
	QueueTypeStandard QueueType = "Standard"
	QueueTypePriority QueueType = "Priority"
)

type QueueMetadata struct {
	Type QueueType
}

type Queue struct {
	Id       string
	Elements CustomerQueue
	Metadata *QueueMetadata
}

type QueueI interface {
	Add(c *customer.Customer) (*customer.Customer, error)
	Pop() (*customer.Customer, error)
}

func NewQueue(ctx context.Context, id string, qtype QueueType) (*Queue, error) {

	return &Queue{
		Id:       id,
		Elements: make(CustomerQueue, 0),
		Metadata: &QueueMetadata{
			Type: qtype,
		},
	}, nil
}

func (q *Queue) Add(c *customer.Customer) {
	currQueue := q.Elements
	currQueue = append(currQueue, c)

	q.Elements = currQueue

}

func (q *Queue) Pop() (c *customer.Customer, err error) {
	currQueue := q.Elements
	len := len(currQueue)
	if len == 0 {
		return nil, fmt.Errorf("cannot pop from empty queue")
	}

	customer := currQueue[0]
	q.Elements = currQueue[1:]

	return customer, nil
}

func (q *Queue) Len() int {
	return len(q.Elements)
}
