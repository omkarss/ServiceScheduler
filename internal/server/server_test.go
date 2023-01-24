package server_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	"github.com/omkar.sunthankar/servicescheduler/internal/server"
	mock_queue "github.com/omkar.sunthankar/servicescheduler/mocks/queue"
	mock_scheduler "github.com/omkar.sunthankar/servicescheduler/mocks/scheduler"
	"github.com/stretchr/testify/suite"
)

type TestServerSuite struct {
	suite.Suite

	ctrl                     *gomock.Controller
	mock_vip_first_scheduler *mock_scheduler.MockSchedulerI
	mock_custom_scheduler    *mock_scheduler.MockSchedulerI
	mock_queue               *mock_queue.MockQueueI
	server                   *server.Server
}

func getStandardCustomer(name string, phoneNumber string, cType customer.CustomerType, ticketNumber int) *customer.StandardCustomer {

	return &customer.StandardCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: ticketNumber,
				Type:         customer.CustomerTypeStandard,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func getVIPCustomer(name string, phoneNumber string, cType customer.CustomerType, tNumber int) *customer.VIPCustomer {

	return &customer.VIPCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: tNumber,
				Type:         customer.CustomerTypeVIP,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func InitServer() *server.Server {
	ctx := context.Background()

	queueS, err := queue.NewQueue(ctx, uuid.NewString(), queue.QueueTypeStandard)
	if err != nil {
		log.Fatal("cannot create standard queue")

	}
	queueP, err := queue.NewQueue(ctx, uuid.NewString(), queue.QueueTypePriority)
	if err != nil {
		log.Fatal("cannot create priority queue")

	}

	// create a VIP Scheduler
	schedulerQueues := make(map[queue.QueueType]*queue.Queue, 0)
	schedulerQueues[queue.QueueTypeStandard] = queueS
	schedulerQueues[queue.QueueTypePriority] = queueP

	SchedulerTypeA, err := scheduler.NewSchedulerTypeA(1)
	if err != nil {
		log.Fatal("Cannot create VIP scheduler")
	}

	// create a Custom Scheduler
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

	SchedulerTypeB, err := scheduler.NewSchedulerTypeB(1, schedulerMetadata)
	if err != nil {
		log.Fatal("Cannot create custom scheduler")
	}

	q := make(map[queue.QueueType]*queue.Queue)
	q[queue.QueueTypeStandard] = queueS
	q[queue.QueueTypePriority] = queueP

	s := server.NewServer(
		SchedulerTypeA,
		SchedulerTypeB,
		q,
	)
	return s
}

func (s *TestServerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.server = InitServer()
	s.mock_queue = mock_queue.NewMockQueueI(s.ctrl)
	s.mock_vip_first_scheduler = mock_scheduler.NewMockSchedulerI(s.ctrl)
	s.mock_custom_scheduler = mock_scheduler.NewMockSchedulerI(s.ctrl)
}

func (s *TestServerSuite) AfterTest() {}

func (s *TestServerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestServerHanlderSuite(t *testing.T) {
	suite.Run(t, new(TestServerSuite))
}
