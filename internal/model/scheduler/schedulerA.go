package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type SchedulerTypeA struct {
	Id           int
	TicketNumber int
	mutex        *sync.Mutex
}

func NewSchedulerTypeA(id int) (*SchedulerTypeA, error) {

	return &SchedulerTypeA{
		Id:    id,
		mutex: &sync.Mutex{},
	}, nil
}

func (sc *SchedulerTypeA) GetNextTicketNumber(ctx context.Context) int {

	sc.TicketNumber++
	return sc.TicketNumber
}

/* Create a lock on scheduler so that only customer is scheduled at 1 time.

Needed for a scenario when multiple workers sharing the same scheduler object .
1. Customer A of standard type is scheduled by worker A .
2. Customer B of priority is checked In and scheduled by worker B .

*/
func (sc *SchedulerTypeA) GetNextCustomer(ctx context.Context, queueMap map[queue.QueueType]*queue.Queue) (*customer.Customer, error) {

	var (
		c   *customer.Customer
		err error
	)

	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Try to get from VIP queue , if no one present get from standard queue.
	sq := queueMap[queue.QueueTypeStandard]
	pq := queueMap[queue.QueueTypePriority]

	if len(pq.Elements) > 0 || len(sq.Elements) > 0 {
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
		}
	} else {
		return nil, fmt.Errorf("no customer available to attend")
	}
	return c, nil

}
