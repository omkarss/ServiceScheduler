package scheduler

import (
	"context"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

type SchedulerQueues map[queue.QueueType]*queue.Queue

type SchedulerI interface {
	GetNextCustomer(*context.Context) (interface{}, error)
	GetNextTicketNumber(*context.Context) int
}
