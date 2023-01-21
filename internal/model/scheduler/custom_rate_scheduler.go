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
	Id       int
	Queues   SchedulerQueues
	Metadata *SchedulerMetadata
}

func NewCustomScheduler(id int, q SchedulerQueues, metadata *SchedulerMetadata) (*CustomScheduler, error) {
	return &CustomScheduler{
		Id:       id,
		Queues:   q,
		Metadata: metadata,
	}, nil
}

func (sc *CustomScheduler) CheckInCustomer(ctx context.Context, c interface{}) (*CustomScheduler, error) {

	switch c.(type) {
	case customer.StandardCustomer:
		sCustomer := c.(customer.StandardCustomer)
		// Check if this key exist in map
		// TODO: always adding to adding first queue, but we can extend this multiple queues
		sq := sc.Queues[queue.QueueTypeStandard]

		sq.Add(&sCustomer.Customer)
		sc.Queues[queue.QueueTypeStandard] = sq

	case customer.VIPCustomer:
		sCustomer := c.(customer.VIPCustomer)
		// Check if this key exist in map
		// TODO: always adding to adding first queue, but we can extend this multiple queues
		vq := sc.Queues[queue.QueueTypePriority]

		vq.Add(&sCustomer.Customer)
		sc.Queues[queue.QueueTypePriority] = vq

	default:
		return nil, fmt.Errorf("customer type is unknown")
	}

	return sc, nil
}

func indexOf(queueType queue.QueueType) int {

	for index, q := range customerAttendingOrder {
		if queueType == q {
			return index
		}
	}

	return -1
}
func (sc *CustomScheduler) AddQueueToScheduler(ctx *context.Context, q interface{}) (*CustomScheduler, error) {
	// TODO : To add a queue to scheduler : Should be handled by admin
	return sc, nil
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
		nextQueueIndex := (index + 1) % len(polledQ)
		sc.Metadata.ShouldPollFromQueue[customerAttendingOrder[nextQueueIndex]] = true
	}

	return nil
}

func (sc *CustomScheduler) GetNextCustomer(context.Context) (interface{}, error) {
	// there should be only 1 queue to poll from
	// Pop From Queue
	// Update the Queue Metadata
	queueToPollFrom := getQueueToPollFrom(sc)
	if queueToPollFrom == "" {
		return nil, fmt.Errorf("issue while polling from queue")
	}

	var (
		customer *customer.Customer
		err      error
	)

	queue, ok := sc.Queues[queueToPollFrom]
	if !ok {
		return nil, fmt.Errorf("cannot get queue")
	}

	customer, err = queue.Pop()
	if err != nil {
		return nil, err
	}

	// Set the count Metadata to len of elements
	sc.Metadata.CurrentPollRemain[queueToPollFrom] = len(queue.Elements)

	sc.updateScheduler(queueToPollFrom)
	return customer, nil
}

func (sc *CustomScheduler) UpdateRate(ctx context.Context) int {
	return 0
}
