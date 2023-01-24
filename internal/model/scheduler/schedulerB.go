package scheduler

import (
	"context"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type SchedulerMetadata struct {
	CurrentPollRemain   map[queue.QueueType]int
	QueuePollRate       map[queue.QueueType]int
	ShouldPollFromQueue map[queue.QueueType]bool
}

type SchedulerTypeB struct {
	Id           int
	TicketNumber int
	Metadata     *SchedulerMetadata
}

func NewSchedulerTypeB(id int, metadata *SchedulerMetadata) (*SchedulerTypeB, error) {
	return &SchedulerTypeB{
		Id:       id,
		Metadata: metadata,
	}, nil
}

func (sc *SchedulerTypeB) AddQueueToScheduler(ctx context.Context, c interface{}) (*SchedulerTypeB, error) {

	return nil, nil
}

func indexOf(queueType queue.QueueType) int {

	for index, q := range customerAttendingOrder {
		if queueType == q {
			return index
		}
	}

	return -1
}

func getQueueToPollFrom(sc *SchedulerTypeB) queue.QueueType {

	for _, qType := range customerAttendingOrder {
		if sc.Metadata.ShouldPollFromQueue[qType] {
			return qType
		}
	}

	return ""
}

func (sc *SchedulerTypeB) updateScheduler(polledQ queue.QueueType) error {

	index := indexOf(polledQ)
	if index == -1 {
		return fmt.Errorf("cannot find index of Queue")
	}
	if sc.Metadata.CurrentPollRemain[polledQ] == 0 {
		// Update the count to full  and set the polling to false
		sc.Metadata.CurrentPollRemain[polledQ] = sc.Metadata.QueuePollRate[polledQ]
		sc.Metadata.ShouldPollFromQueue[polledQ] = false

		// Update the nextQueue to true so that we can poll from it
		nextQueueIndex := (index + 1) % len(customerAttendingOrder)
		sc.Metadata.ShouldPollFromQueue[customerAttendingOrder[nextQueueIndex]] = true
	}

	return nil
}

func (sc *SchedulerTypeB) GetNextCustomer(ctx context.Context, queueMap map[queue.QueueType]*queue.Queue) (*customer.Customer, error) {

	var (
		customer *customer.Customer
		err      error
	)

	queueToPollFrom := getQueueToPollFrom(sc)
	if queueToPollFrom == "" {
		return nil, fmt.Errorf("issue while polling from queue")
	}

	queue := queueMap[queueToPollFrom]

	// if current queue is empty get from next queue ( doesnt make sense to make standard customers waiting )
	if len(queue.Elements) == 0 {

		// Poll from next available queue according to priority
		i := indexOf(queueToPollFrom)
		next_queue := (i + 1) % len(customerAttendingOrder)

		for {
			if next_queue == i {
				return nil, fmt.Errorf("no customer to attend")
			}

			qType := customerAttendingOrder[next_queue]

			if len(queueMap[qType].Elements) > 0 {
				queue = queueMap[qType]
				break
			}

			next_queue = (next_queue + 1) % len(customerAttendingOrder)
		}
	}

	customer, err = queue.Pop()
	if err != nil {
		return nil, err
	}

	// Decrease the poll by 1
	sc.Metadata.CurrentPollRemain[queueToPollFrom] -= 1

	sc.updateScheduler(queueToPollFrom)
	return customer, nil
}

func (sc *SchedulerTypeB) UpdateRate(ctx context.Context) int {
	return 0
}
