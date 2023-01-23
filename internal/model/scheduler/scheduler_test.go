package scheduler_test

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	mock_queue "github.com/omkar.sunthankar/servicescheduler/mocks/queue"
	"github.com/stretchr/testify/suite"
)

type SchedulerSuite struct {
	suite.Suite

	ctrl       *gomock.Controller
	mock_queue *mock_queue.MockQueueI
	queues     map[queue.QueueType]*queue.Queue
	schedulerA *scheduler.VIPFirstSceduler
	schedulerB *scheduler.CustomScheduler
}

func (s *SchedulerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mock_queue = mock_queue.NewMockQueueI(s.ctrl)
	s.queues = make(map[queue.QueueType]*queue.Queue, 0)
	s.schedulerA, _ = scheduler.NewVIPFirstSceduler(1)
	s.schedulerB, _ = scheduler.NewCustomScheduler(1, nil)
}

func (s *SchedulerSuite) SetupIndividualTest() {
	s.queues = make(map[queue.QueueType]*queue.Queue, 0)
}

func (s *SchedulerSuite) AfterTest() {
	s.queues = nil
}

func getQueue(qT queue.QueueType) *queue.Queue {
	queueS, err := queue.NewQueue(context.Background(), uuid.NewString(), qT)
	if err != nil {
		log.Fatal("cannot create standard queue")

	}

	return queueS
}

func (s *SchedulerSuite) BeforeIndividualTest() {

	queueS := getQueue(queue.QueueTypeStandard)
	queueP := getQueue(queue.QueueTypePriority)

	s.queues[queue.QueueTypeStandard] = queueS
	s.queues[queue.QueueTypePriority] = queueP
}

func (s *SchedulerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSchedulerSuite(t *testing.T) {
	suite.Run(t, new(SchedulerSuite))
}
