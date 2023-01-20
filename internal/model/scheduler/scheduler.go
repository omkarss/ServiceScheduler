package scheduler

import (
	"context"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type SchedulerQueues map[queue.QueueType]interface{}

type SchedulerI interface {
	CheckInCustomer(*context.Context, interface{}) (interface{}, error)
	AddQueueToScheduler(*context.Context, interface{}) (interface{}, error)
	GetNextCustomer(*context.Context) (interface{}, error)
}
