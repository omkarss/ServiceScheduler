package queue

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

type VipFirstQueue struct {
	Id       string
	Elements VIPCustomers
}

func NewVipFirstQueue(ctx context.Context, id string) (*VipFirstQueue, error) {

	return &VipFirstQueue{
		Id:       id,
		Elements: make(VIPCustomers, 0),
	}, nil
}

func (q *VipFirstQueue) Add(c *customer.VIPCustomer) {
	currQueue := q.Elements
	currQueue = append(currQueue, c)

	q.Elements = currQueue

}

func (q *VipFirstQueue) Pop() (*customer.VIPCustomer, error) {
	currQueue := q.Elements
	len := len(currQueue)
	if len == 0 {
		return nil, fmt.Errorf("cannot pop from empty queue")
	}

	customer := currQueue[len-1]
	q.Elements = currQueue[:(len - 1)]

	return customer, nil
}

func (q *VipFirstQueue) Len() int {
	return len(q.Elements)
}
