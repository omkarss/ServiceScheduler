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

func NewVIPFirstSceduler(id int, q SchedulerQueues) (*VIPFirstSceduler, error) {

	return &VIPFirstSceduler{
		Id:     id,
		Queues: q,
	}, nil
}

func (sc *VIPFirstSceduler) CheckInCustomer(ctx context.Context, c interface{}) (*VIPFirstSceduler, error) {

	switch c.(type) {
	case customer.StandardCustomer:
		sCustomer := c.(customer.StandardCustomer)
		// TODO: Check if this key exist in map
		sq := sc.Queues[queue.QueueTypeStandard]

		sq.Add(&sCustomer.Customer)
		sc.Queues[queue.QueueTypeStandard] = sq

	case customer.VIPCustomer:
		sCustomer := c.(customer.VIPCustomer)
		// TODO: Check if this key exist in map
		vq := sc.Queues[queue.QueueTypePriority]

		vq.Add(&sCustomer.Customer)
		sc.Queues[queue.QueueTypePriority] = vq

	default:
		return nil, fmt.Errorf("customer type is unknown")
	}

	return sc, nil
}

func (sc *VIPFirstSceduler) AddQueueToScheduler(ctx context.Context, q queue.Queue) error {
	return nil
}

func (p *VIPFirstSceduler) GetNextCustomer(context.Context) (interface{}, error) {

	// check if VIP QUeue has anything in it If yes remove from VIP queue else from Standard Queue
	var (
		customer interface{}
		err      error
	)

	// TODO: Make this a switch if type of queues increase
	vipQueue, ok := p.Queues[queue.QueueTypePriority]
	if !ok {
		return nil, fmt.Errorf("cannot get queue")
	}

	stdQueue, ok := p.Queues[queue.QueueTypeStandard]
	if !ok {
		return nil, fmt.Errorf("cannot get queue")
	}

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
