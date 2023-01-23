package scheduler

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type VIPFirstSceduler struct {
	Id           int
	TicketNumber int
}

func NewVIPFirstSceduler(id int) (*VIPFirstSceduler, error) {

	return &VIPFirstSceduler{
		Id: id,
	}, nil
}

func (sc *VIPFirstSceduler) GetNextTicketNumber(ctx context.Context) int {

	sc.TicketNumber = sc.TicketNumber + 1
	return sc.TicketNumber
}

func (sc *VIPFirstSceduler) GetNextCustomer(ctx context.Context, queueMap map[queue.QueueType]*queue.Queue) (*customer.Customer, error) {

	var c *customer.Customer
	var err error

	// Try to get from VIP queue , if no one present get from standard queue.
	sq := queueMap[queue.QueueTypeStandard]
	pq := queueMap[queue.QueueTypePriority]

	if len(pq.Elements) > 0 {

		c, err = pq.Pop()
		if err != nil {
			return nil, err
		}
	} else if len(sq.Elements) > 0 {
		c, err = sq.Pop()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no customer available to attend")
	}
	return c, nil
}
