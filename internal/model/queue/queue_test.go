package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/stretchr/testify/assert"
)

func getNewQueue(id string, qT queue.QueueType) *queue.Queue {

	q, _ := queue.NewQueue(context.Background(), id, qT)
	return q
}
func TestAddToQueue(t *testing.T) {

	c1 := &customer.Customer{
		FullName:    "Omkar",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeStandard,
			EntryTime:    time.Now(),
		},
	}

	t.Run("Succesfully adds an element to Queue", func(t *testing.T) {
		q := getNewQueue(uuid.NewString(), queue.QueueTypeStandard)

		q.Add(c1)

		// Assertions
		assert.Equal(t, len(q.Elements), 1)
		assert.Equal(t, q.Elements[0].FullName, "Omkar")
		assert.Equal(t, q.Elements[0].PhoneNumber, "xxxx")

	})

}

func TestExists(t *testing.T) {

	c1 := &customer.Customer{
		FullName:    "Omkar",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeStandard,
			EntryTime:    time.Now(),
		},
	}

	t.Run("Returns True for already existing customer", func(t *testing.T) {
		q := getNewQueue(uuid.NewString(), queue.QueueTypeStandard)

		q.Add(c1)

		ex := q.Exists(c1.FullName, c1.PhoneNumber)
		// Assertions
		assert.Equal(t, len(q.Elements), 1)
		assert.Equal(t, q.Elements[0].FullName, "Omkar")
		assert.Equal(t, q.Elements[0].PhoneNumber, "xxxx")
		assert.Equal(t, ex, true)

	})

	t.Run("Returns false for customer which is not present", func(t *testing.T) {
		q := getNewQueue(uuid.NewString(), queue.QueueTypeStandard)

		q.Add(c1)

		ex := q.Exists("B", c1.PhoneNumber)
		// Assertions
		assert.Equal(t, len(q.Elements), 1)
		assert.Equal(t, q.Elements[0].FullName, "Omkar")
		assert.Equal(t, q.Elements[0].PhoneNumber, "xxxx")
		assert.Equal(t, ex, false)

	})

}

func TestPopFromQueue(t *testing.T) {

	c1 := &customer.Customer{
		FullName:    "Omkar",
		PhoneNumber: "xxxx",
		Metadata: &customer.CustomerMetadata{
			TicketNumber: 1,
			Type:         customer.CustomerTypeStandard,
			EntryTime:    time.Now(),
		},
	}
	t.Run("Sucessfully Pops from Queue", func(t *testing.T) {

		q := getNewQueue(uuid.NewString(), queue.QueueTypeStandard)

		q.Add(c1)
		q1 := *q
		assert.Equal(t, len(q.Elements), 1)

		expectedC, _ := q1.Pop()

		// Assertions
		assert.Equal(t, 0, len(q1.Elements))
		assert.Equal(t, expectedC.FullName, "Omkar")
		assert.Equal(t, expectedC.PhoneNumber, "xxxx")

	})

	t.Run("Fails when tries to Pop from Empty Queue", func(t *testing.T) {

		q := getNewQueue(uuid.NewString(), queue.QueueTypeStandard)

		assert.Equal(t, len(q.Elements), 0)

		_, err := q.Pop()

		// Assertions
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot pop from empty queue")

	})

}
