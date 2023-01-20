package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	mock_queue "github.com/omkar.sunthankar/servicescheduler/mocks/queue"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {

	// Create a Standard Customer
	standardCustomer := &customer.StandardCustomer{
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

	t.Run("Adds a Queue to Scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		_ = mock_queue.NewMockQueue(ctrl)

		ctx := context.Background()

		// Create a Standard Queue and Priority Queues
		stdQ, _ := queue.NewStandardQueue(ctx, "1")
		vipQ, _ := queue.NewVipFirstQueue(ctx, "1")

		// Add Customers to Queues
		stdQ.Add(standardCustomer)
		vipQ.Add(vipCustomer)

		// Create a Scheduler
		p, err := scheduler.NewVIPFirstSceduler(1)
		if err != nil {
			t.Error("error occured while creating a scheduler", err)
		}

		err = p.AddQueueToScheduler(ctx, stdQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		err = p.AddQueueToScheduler(ctx, vipQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		assert.NotNil(p)

		sq := p.Queues[queue.QueueTypeStandard].(*queue.StandardQueue)
		vq := p.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)

		assert.Equal(stdQ, sq)
		assert.Equal(vipQ, vq)

	})

	t.Run("Adds a Customer to priority Scheduler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		_ = mock_queue.NewMockQueue(ctrl)

		ctx := context.Background()

		// Create a Standard Queue and Priority Queues
		stdQ, _ := queue.NewStandardQueue(ctx, "1")
		vipQ, _ := queue.NewVipFirstQueue(ctx, "1")

		// Create a Scheduler
		p, err := scheduler.NewVIPFirstSceduler(1)
		if err != nil {
			t.Error("error occured while creating a scheduler", err)
		}

		// Add queue to scheduler
		err = p.AddQueueToScheduler(ctx, stdQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		err = p.AddQueueToScheduler(ctx, vipQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		p, err = p.CheckInCustomer(ctx, *standardCustomer)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		p, err = p.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		expectedstdQueue := p.Queues[queue.QueueTypeStandard].(*queue.StandardQueue)
		expectedprioQueue := p.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)

		assert.Contains(expectedstdQueue.Elements[0].Customer.FullName, "Bob")
		assert.Contains(expectedprioQueue.Elements[0].Customer.FullName, "Alice")
	})

	t.Run("Gets next Customer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		assert := assert.New(t)
		defer ctrl.Finish()

		_ = mock_queue.NewMockQueue(ctrl)

		ctx := context.Background()

		// Create a Standard Queue and Priority Queues
		stdQ, _ := queue.NewStandardQueue(ctx, "1")
		vipQ, _ := queue.NewVipFirstQueue(ctx, "1")

		// Create a Scheduler
		p, err := scheduler.NewVIPFirstSceduler(1)
		if err != nil {
			t.Error("error occured while creating a scheduler", err)
		}

		// Add queue to scheduler
		err = p.AddQueueToScheduler(ctx, stdQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		err = p.AddQueueToScheduler(ctx, vipQ)
		if err != nil {
			t.Error("error occured while adding queue to scheduler", err)
		}

		p, err = p.CheckInCustomer(ctx, *standardCustomer)
		if err != nil {
			t.Error("error occured while checking in a standard customer", err)
		}

		// check this
		p, err = p.CheckInCustomer(ctx, *vipCustomer)
		if err != nil {
			t.Error("error occured while checking in a vip customer", err)
		}

		c, err := p.GetNextCustomer(ctx)
		if err != nil {
			t.Error("error occured while popping from queue", err)
		}
		expectedPQ := p.Queues[queue.QueueTypePriority].(*queue.VipFirstQueue)

		assert.Contains(c.(*customer.VIPCustomer).Customer.FullName, "Alice A")
		assert.Equal(len(expectedPQ.Elements), 0)
	})
}
