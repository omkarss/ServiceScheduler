package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

func (s *SchedulerSuite) TestSchedulerBGetNextCustomer() {

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
			Type:         customer.CustomerTypeStandard,
			EntryTime:    time.Now(),
		},
	}

	c3 := &customer.Customer{
		FullName:    "C",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	c4 := &customer.Customer{
		FullName:    "D",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	c5 := &customer.Customer{
		FullName:    "E",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	c6 := &customer.Customer{
		FullName:    "F",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	c7 := &customer.Customer{
		FullName:    "G",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeVIP,
			EntryTime:    time.Now(),
		},
	}

	defer s.ctrl.Finish()

	s.T().Run("Gets a customer from  VIP queue first ", func(t *testing.T) {
		s.SetupIndividualTest()
		s.BeforeIndividualTest()
		defer s.ctrl.Finish()

		s.queues[queue.QueueTypeStandard].Add(c1)
		s.queues[queue.QueueTypeStandard].Add(c2)
		s.queues[queue.QueueTypePriority].Add(c3)
		s.queues[queue.QueueTypePriority].Add(c4)
		s.queues[queue.QueueTypePriority].Add(c5)
		s.queues[queue.QueueTypePriority].Add(c6)
		s.queues[queue.QueueTypePriority].Add(c7)

		expectedC, err := s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c3.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c4.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 2)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], false)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], true)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c1.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 2)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c5.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c6.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 2)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], false)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], true)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c2.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 2)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)

		expectedC, err = s.schedulerB.GetNextCustomer(context.Background(), s.queues)
		r.NoError(err)
		r.Equal(c7.FullName, expectedC.FullName)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypePriority], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypePriority], true)
		r.Equal(s.schedulerB.Metadata.CurrentPollRemain[queue.QueueTypeStandard], 1)
		r.Equal(s.schedulerB.Metadata.ShouldPollFromQueue[queue.QueueTypeStandard], false)

	})

	s.T().Run("Gets a customer from standard queue when vip queue has no customers", func(t *testing.T) {
		s.SetupIndividualTest()
		s.BeforeIndividualTest()
		defer s.ctrl.Finish()

		s.queues[queue.QueueTypeStandard].Add(c1)

		expectedC, err := s.schedulerB.GetNextCustomer(context.Background(), s.queues)

		r.NoError(err)
		r.Equal(c1.FullName, expectedC.FullName)
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
