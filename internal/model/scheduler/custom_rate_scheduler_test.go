package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	"github.com/stretchr/testify/assert"
)

// Create a Standard Customer
func getStandardCustomer(name string, phoneNumber string, cType customer.CustomerType) *customer.StandardCustomer {

	return &customer.StandardCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeStandard,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func getVIPCustomer(name string, phoneNumber string, cType customer.CustomerType) *customer.VIPCustomer {

	return &customer.VIPCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeVIP,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func getQueue(queueType queue.QueueType) *queue.Queue {
	return &queue.Queue{
		Id:       uuid.NewString(),
		Elements: make(queue.CustomerQueue, 0),
		Metadata: &queue.QueueMetadata{
			Type: queueType,
		},
	}
}

func getCustomScheduler(queues []*queue.Queue) (*scheduler.CustomScheduler, error) {

	// Create a Scheduler
	schedulerQueues := make(map[queue.QueueType]*queue.Queue, 0)

	for _, q := range queues {
		schedulerQueues[q.Metadata.Type] = q
	}

	schedulerMetadata := &scheduler.SchedulerMetadata{}

	schedulerMetadata.ShouldPollFromQueue = make(map[queue.QueueType]bool, 0)
	schedulerMetadata.CurrentPollRemain = make(map[queue.QueueType]int, 0)
	schedulerMetadata.QueuePollRate = make(map[queue.QueueType]int, 0)

	schedulerMetadata.CurrentPollRemain[queue.QueueTypeStandard] = 1
	schedulerMetadata.QueuePollRate[queue.QueueTypeStandard] = 1
	schedulerMetadata.ShouldPollFromQueue[queue.QueueTypeStandard] = false

	schedulerMetadata.QueuePollRate[queue.QueueTypePriority] = 2
	schedulerMetadata.CurrentPollRemain[queue.QueueTypePriority] = 2
	schedulerMetadata.ShouldPollFromQueue[queue.QueueTypePriority] = true

	p, err := scheduler.NewCustomScheduler(1, schedulerQueues, schedulerMetadata)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func TestCustomRateSchedulerCheckInCustomer(t *testing.T) {

	standardCustomer1 := getStandardCustomer("Bob", "1234", customer.CustomerTypeStandard)
	standardCustomer2 := getStandardCustomer("Cathy", "1234", customer.CustomerTypeStandard)
	vipCustomer := getVIPCustomer("Alice", "1234", customer.CustomerTypeVIP)

	standardQ := getQueue(queue.QueueTypeStandard)
	priorityQ := getQueue(queue.QueueTypePriority)

	t.Run("Adds a customer to priority scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		ctx := context.Background()
		customScheduler, err := getCustomScheduler([]*queue.Queue{standardQ, priorityQ})

		// asserts to make sure we have the right scheduler
		assert.Nil(err)
		assert.NotNil(customScheduler)
		assert.Equal(len(customScheduler.Queues), 2)

		customScheduler, err = customScheduler.CheckInCustomer(ctx, *standardCustomer1)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		customScheduler, err = customScheduler.CheckInCustomer(ctx, *standardCustomer2)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		customScheduler, err = customScheduler.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		assert.Equal(len(customScheduler.Queues), 2)
		assert.Contains(customScheduler.Queues[queue.QueueTypeStandard].Elements[0].FullName, "Bob")
		assert.Contains(customScheduler.Queues[queue.QueueTypeStandard].Elements[1].FullName, "Cathy")
		assert.Contains(customScheduler.Queues[queue.QueueTypePriority].Elements[0].FullName, "Alice")
	})

}

func TestCustomRateSchedulerGetNextCustomer(t *testing.T) {

	standardCustomer1 := getStandardCustomer("Bob", "1234", customer.CustomerTypeStandard)
	standardCustomer2 := getStandardCustomer("Cathy", "1234", customer.CustomerTypeStandard)
	vipCustomer := getVIPCustomer("Alice", "1234", customer.CustomerTypeVIP)
	vipCustomer2 := getVIPCustomer("VIPAlice", "1234", customer.CustomerTypeVIP)

	standardQ := getQueue(queue.QueueTypeStandard)
	priorityQ := getQueue(queue.QueueTypePriority)

	t.Run("gets the next customer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		ctx := context.Background()

		customScheduler, err := getCustomScheduler([]*queue.Queue{standardQ, priorityQ})

		// asserts to make sure we have the right scheduler
		assert.Nil(err)
		assert.NotNil(customScheduler)
		assert.Equal(len(customScheduler.Queues), 2)

		customScheduler, err = customScheduler.CheckInCustomer(ctx, *standardCustomer1)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		customScheduler, err = customScheduler.CheckInCustomer(ctx, *standardCustomer2)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		customScheduler, err = customScheduler.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		customScheduler, err = customScheduler.CheckInCustomer(ctx, *vipCustomer2)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		c1, err := customScheduler.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}

		c2, err := customScheduler.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}

		c3, err := customScheduler.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}
		assert.Contains(c1.(*customer.Customer).FullName, "Alice")
		assert.Contains(c2.(*customer.Customer).FullName, "VIPAlice")
		assert.Contains(c3.(*customer.Customer).FullName, "Bob")

		assert.Equal(len(customScheduler.Queues[queue.QueueTypePriority].Elements), 0)
		assert.Equal(len(customScheduler.Queues[queue.QueueTypeStandard].Elements), 1)
	})
}

func TestCustomRateSchedulerAdd(t *testing.T) {

	standardQ := getQueue(queue.QueueTypeStandard)
	priorityQ := getQueue(queue.QueueTypePriority)

	t.Run("Succesfully Initializes Scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		schedulerMetadata := &scheduler.SchedulerMetadata{}

		// Create a Scheduler
		schedulerQueues := make(map[queue.QueueType]*queue.Queue, 0)
		schedulerQueues[queue.QueueTypeStandard] = standardQ
		schedulerQueues[queue.QueueTypePriority] = priorityQ

		schedulerMetadata.ShouldPollFromQueue = make(map[queue.QueueType]bool, 0)
		schedulerMetadata.CurrentPollRemain = make(map[queue.QueueType]int, 0)
		schedulerMetadata.QueuePollRate = make(map[queue.QueueType]int, 0)

		schedulerMetadata.CurrentPollRemain[queue.QueueTypeStandard] = 1
		schedulerMetadata.QueuePollRate[queue.QueueTypeStandard] = 1
		schedulerMetadata.ShouldPollFromQueue[queue.QueueTypeStandard] = false

		schedulerMetadata.QueuePollRate[queue.QueueTypePriority] = 2
		schedulerMetadata.CurrentPollRemain[queue.QueueTypePriority] = 2
		schedulerMetadata.ShouldPollFromQueue[queue.QueueTypePriority] = true

		scheduler, err := scheduler.NewCustomScheduler(1, schedulerQueues, schedulerMetadata)

		assert.NoError(err)
		assert.NotNil(scheduler)
		assert.Equal(scheduler.Id, 1)

		assert.Equal(scheduler.Metadata.QueuePollRate[queue.QueueTypePriority], 2)
		assert.Equal(scheduler.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		assert.Equal(scheduler.Metadata.CurrentPollRemain[queue.QueueTypePriority], 2)

		assert.Equal(scheduler.Metadata.QueuePollRate[queue.QueueTypeStandard], 1)
		assert.Equal(scheduler.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)
		assert.Equal(scheduler.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
	})

}
