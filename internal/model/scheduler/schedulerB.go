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

type CustomScheduler struct {
	Id           int
	TicketNumber int
	Metadata     *SchedulerMetadata
}

func NewCustomScheduler(id int, metadata *SchedulerMetadata) (*CustomScheduler, error) {
	return &CustomScheduler{
		Id:       id,
		Metadata: metadata,
	}, nil
}

func (sc *CustomScheduler) AddQueueToScheduler(ctx context.Context, c interface{}) (*CustomScheduler, error) {

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

func getQueueToPollFrom(sc *CustomScheduler) queue.QueueType {

	for _, qType := range customerAttendingOrder {
		if sc.Metadata.ShouldPollFromQueue[qType] {
			return qType
		}
	}

	return ""
}

func (sc *CustomScheduler) updateScheduler(polledQ queue.QueueType) error {

	index := indexOf(polledQ)
	if index == -1 {
		return fmt.Errorf("cannot find index of Queue")
	}
	if sc.Metadata.CurrentPollRemain[polledQ] == 0 {
		// Update the count to full by 1 and set the polling to false
		sc.Metadata.CurrentPollRemain[polledQ] = sc.Metadata.QueuePollRate[polledQ]
		sc.Metadata.ShouldPollFromQueue[polledQ] = false

		// Update the nextQueue to true so that we can poll from it
		nextQueueIndex := (index + 1) % len(customerAttendingOrder)
		sc.Metadata.ShouldPollFromQueue[customerAttendingOrder[nextQueueIndex]] = true
	}

	return nil
}

func (sc *CustomScheduler) GetNextCustomer(ctx context.Context, queueMap map[queue.QueueType]*queue.Queue) (interface{}, error) {
	// there should be only 1 queue to poll from
	// Pop From Queue
	// Update the Queue Metadata
	var (
		customer *customer.Customer
		err      error
	)

	queueToPollFrom := getQueueToPollFrom(sc)
	if queueToPollFrom == "" {
		return nil, fmt.Errorf("issue while polling from queue")
	}

	queue := queueMap[queueToPollFrom]
	customer, err = queue.Pop()
	if err != nil {
		return nil, err
	}

	// Decrease the poll by 1
	sc.Metadata.CurrentPollRemain[queueToPollFrom] -= 1

	sc.updateScheduler(queueToPollFrom)
	return customer, nil
}

func (sc *CustomScheduler) UpdateRate(ctx context.Context) int {
	return 0
}
