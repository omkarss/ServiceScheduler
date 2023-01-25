package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

func (s *SchedulerSuite) TestSchedulerAGetNextCustomer() {

	r := s.Require()

	c1 := &customer.Customer{
		FullName:    "A",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeStandard,
			EntryTime:    time.Now(),
		},
	}

	c2 := &customer.Customer{
		FullName:    "B",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	_ = &customer.Customer{
		FullName:    "C",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}
	defer s.ctrl.Finish()

	s.T().Run("Gets a customer from  priority queue first ", func(t *testing.T) {
		s.SetupIndividualTest()
		s.BeforeIndividualTest()
		defer s.ctrl.Finish()

		s.queues[queue.QueueTypeStandard].Add(c1)
		s.queues[queue.QueueTypePriority].Add(c2)

		// s.mock_queue.EXPECT().Pop().Return(c2, nil)

		expectedC, err := s.schedulerA.GetNextCustomer(context.Background(), s.queues)

		r.NoError(err)
		r.Equal(c2.FullName, expectedC.FullName)
		r.Equal(0, len(s.queues[queue.QueueTypePriority].Elements))
	})

	s.T().Run("Gets a customer from standard when vip queue has no customers", func(t *testing.T) {
		s.SetupIndividualTest()
		s.BeforeIndividualTest()
		defer s.ctrl.Finish()

		s.queues[queue.QueueTypeStandard].Add(c1)

		// s.mock_queue.EXPECT().Pop().Return(c1, nil)

		expectedC, err := s.schedulerA.GetNextCustomer(context.Background(), s.queues)

		r.NoError(err)
		r.Equal(c1.FullName, expectedC.FullName)
		r.Equal(0, len(s.queues[queue.QueueTypeStandard].Elements))
	})

	s.T().Run("does not schedule when there is no customer to poll ", func(t *testing.T) {
		s.SetupIndividualTest()
		s.BeforeIndividualTest()
		defer s.ctrl.Finish()

		_, err := s.schedulerA.GetNextCustomer(context.Background(), s.queues)

		r.Error(err)
		r.EqualError(err, "no customer available to attend")

	})
}
