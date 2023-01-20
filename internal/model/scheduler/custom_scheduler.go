package scheduler

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type CustomScheduler struct {
	Id     int
	Queues SchedulerQueues
}

func NewCustomScheduler(id int) (*VIPFirstSceduler, error) {
	return &VIPFirstSceduler{
		Id:     id,
		Queues: make(SchedulerQueues, 0),
	}, nil
}

func (sc *CustomScheduler) CheckInCustomer(ctx *context.Context, c interface{}) (*CustomScheduler, error) {

	switch c.(type) {
	case customer.StandardCustomer:
		sCustomer := c.(customer.StandardCustomer)
		sq, ok := sc.Queues[queue.QueueTypeStandard].(*queue.StandardQueue)
		if !ok {
			return nil, fmt.Errorf("cannot cast to Standard Queue")
		}

		sq.Add(&sCustomer)
		sc.Queues[queue.QueueTypeStandard] = sCustomer

	case customer.VIPCustomer:
		vipCustomer := c.(customer.VIPCustomer)
		vq, ok := sc.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)
		if !ok {
			return nil, fmt.Errorf("cannot cast to VIP Queue")
		}

		vq.Add(&vipCustomer)
		sc.Queues[queue.QueueTypePriority] = vipCustomer

	default:
		return nil, fmt.Errorf("customer type is unknown")
	}

	return sc, nil

}

func (sc *CustomScheduler) AddQueueToScheduler(ctx *context.Context, q interface{}) (*CustomScheduler, error) {

	switch q.(type) {
	case queue.StandardQueue:
		sQueue := q.(queue.StandardQueue)
		sc.Queues[queue.QueueTypeStandard] = sQueue

	case queue.VipFirstQueue:
		vipQueue := q.(queue.VipFirstQueue)
		sc.Queues[queue.QueueTypePriority] = vipQueue

	default:
		return nil, fmt.Errorf("queue type is unknown")
	}

	return sc, nil
}

func (*CustomScheduler) GetNextCustomer(*context.Context) (*customer.Customer, error) {
	return nil, nil
}
