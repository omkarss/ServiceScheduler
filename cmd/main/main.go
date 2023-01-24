package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	"github.com/omkar.sunthankar/servicescheduler/internal/server"
)

const port = ":8080"

func main() {

	ctx := context.Background()

	//  Create a Queues
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

	// Assign to server struct
	s := server.NewServer(
		SchedulerTypeA,
		SchedulerTypeB,
		q,
	)

	if err := http.ListenAndServe(port, s); err != nil {
		log.Fatal(err)
	}

}
