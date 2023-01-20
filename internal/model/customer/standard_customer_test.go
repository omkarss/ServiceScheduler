package customer_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	mock_queue "github.com/omkar.sunthankar/servicescheduler/mocks/queue"
)

func TestProcess(t *testing.T) {

	t.Run("Creates a Customer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = mock_queue.NewMockQueue(ctrl)

		_ = context.Background()

		// Create a Customer
		_ = &customer.StandardCustomer{
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

		_ = &customer.VIPCustomer{
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

	})
}
