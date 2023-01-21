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

func getVIPScheduler(queues []*queue.Queue) (*scheduler.VIPFirstSceduler, error) {

	// Create a Scheduler
	schedulerQueues := make(map[queue.QueueType]*queue.Queue, 0)

	for _, q := range queues {
		schedulerQueues[q.Metadata.Type] = q
	}

	p, err := scheduler.NewVIPFirstSceduler(1, schedulerQueues)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func TestVIPFirstScheduler(t *testing.T) {

	// Create a Standard Customer
	standardCustomer1 := &customer.StandardCustomer{
		Customer: customer.Customer{
			FullName:    "Bob A",
			PhoneNumber: "123456",
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeStandard,
				EntryTime:    time.Now().UTC(),
			},
		},
	}

	standardCustomer2 := &customer.StandardCustomer{
		Customer: customer.Customer{
			FullName:    "Cathy K",
			PhoneNumber: "123456",
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeStandard,
				EntryTime:    time.Now().UTC(),
			},
		},
	}

	// Create a VIP Customer
	vipCustomer := &customer.VIPCustomer{
		Customer: customer.Customer{
			FullName:    "Alice A",
			PhoneNumber: "123456",
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeVIP,
				EntryTime:    time.Now().UTC(),
			},
		},
	}

	// Queues
	standardQ := &queue.Queue{
		Id:       uuid.NewString(),
		Elements: make(queue.CustomerQueue, 0),
		Metadata: &queue.QueueMetadata{
			Type: queue.QueueTypeStandard,
		},
	}

	priorityQ := &queue.Queue{
		Id:       uuid.NewString(),
		Elements: make(queue.CustomerQueue, 0),
		Metadata: &queue.QueueMetadata{
			Type: queue.QueueTypePriority,
		},
	}

	t.Run("Succesfully Initializes Scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		// Add Customers to Queues
		standardQ.Add(&standardCustomer1.Customer)
		priorityQ.Add(&vipCustomer.Customer)

		// Create a Scheduler
		schedulerQueues := make(map[queue.QueueType]*queue.Queue, 0)
		schedulerQueues[queue.QueueTypeStandard] = standardQ
		schedulerQueues[queue.QueueTypePriority] = priorityQ

		scheduler, err := scheduler.NewVIPFirstSceduler(1, schedulerQueues)

		assert.NoError(err)
		assert.NotNil(scheduler)
		assert.Equal(scheduler.Id, 1)

	})

	t.Run("Adds a customer to priority scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		ctx := context.Background()
		vipScheduler, err := getVIPScheduler([]*queue.Queue{standardQ, priorityQ})

		// asserts to make sure we have the right scheduler
		assert.Nil(err)
		assert.NotNil(vipScheduler)
		assert.Equal(len(vipScheduler.Queues), 2)

		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *standardCustomer1)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *standardCustomer2)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		assert.Equal(len(vipScheduler.Queues), 2)
		assert.Contains(vipScheduler.Queues[queue.QueueTypeStandard].Elements[0].FullName, "Bob")
		assert.Contains(vipScheduler.Queues[queue.QueueTypeStandard].Elements[1].FullName, "Cathy K")
		assert.Contains(vipScheduler.Queues[queue.QueueTypePriority].Elements[0].FullName, "Alice A")
	})

	t.Run("gets the next customer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		ctx := context.Background()

		vipScheduler, err := getVIPScheduler([]*queue.Queue{standardQ, priorityQ})

		// asserts to make sure we have the right scheduler
		assert.Nil(err)
		assert.NotNil(vipScheduler)
		assert.Equal(len(vipScheduler.Queues), 2)

		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *standardCustomer1)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *standardCustomer2)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		vipScheduler, err = vipScheduler.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		c1, err := vipScheduler.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}

		c2, err := vipScheduler.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}

		assert.Contains(c1.(*customer.Customer).FullName, "Alice A")
		assert.Contains(c2.(*customer.Customer).FullName, "Bob A")
		assert.Equal(len(vipScheduler.Queues[queue.QueueTypePriority].Elements), 0)
		assert.Equal(len(vipScheduler.Queues[queue.QueueTypeStandard].Elements), 1)
	})
}
