package scheduler

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type VIPFirstSceduler struct {
	Id     int
	Queues SchedulerQueues
}

func NewVIPFirstSceduler(id int) (*VIPFirstSceduler, error) {
	return &VIPFirstSceduler{
		Id:     id,
		Queues: make(SchedulerQueues, 0),
	}, nil
}

func (sc *VIPFirstSceduler) CheckInCustomer(ctx context.Context, c interface{}) (*VIPFirstSceduler, error) {

	switch c.(type) {
	case customer.StandardCustomer:
		sCustomer := c.(customer.StandardCustomer)
		sq, ok := sc.Queues[queue.QueueTypeStandard].(*queue.StandardQueue)
		if !ok {
			return nil, fmt.Errorf("cannot cast to Standard Queue")
		}

		sq.Add(&sCustomer)
		sc.Queues[queue.QueueTypeStandard] = sq

	case customer.VIPCustomer:
		vipCustomer := c.(customer.VIPCustomer)
		vq, ok := sc.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)
		if !ok {
			return nil, fmt.Errorf("cannot cast to VIP Queue")
		}

		vq.Add(&vipCustomer)
		sc.Queues[queue.QueueTypePriority] = vq

	default:
		return nil, fmt.Errorf("customer type is unknown")
	}

	return sc, nil
}

func (sc *VIPFirstSceduler) AddQueueToScheduler(ctx context.Context, q interface{}) error {

	switch q.(type) {
	case *queue.StandardQueue:
		sQueue := q.(*queue.StandardQueue)
		sc.Queues[queue.QueueTypeStandard] = sQueue

	case *queue.VipFirstQueue:
		vipQueue := q.(*queue.VipFirstQueue)
		sc.Queues[queue.QueueTypePriority] = vipQueue

	default:
		return fmt.Errorf("queue type is unknown")
	}

	return nil
}

func (p *VIPFirstSceduler) GetNextCustomer(context.Context) (interface{}, error) {
	// check if VIP QUeue has anything in it If yes remove from VIP queue else from Standard Queue

	vipQueue, ok := p.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)
	if !ok {
		return nil, fmt.Errorf("cannot get queue")
	}
	stdQueue, ok := p.Queues[queue.QueueTypeStandard].(*queue.StandardQueue)
	if !ok {
		return nil, fmt.Errorf("cannot get queue")
	}

	var (
		customer interface{}
		err      error
	)

	if len(vipQueue.Elements) > 0 {
		customer, err = vipQueue.Pop()
		if err != nil {
			return nil, err
		}

	} else if len(stdQueue.Elements) > 0 {
		customer, err = stdQueue.Pop()
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("no customers to attend")
	}

	return customer, nil
}
