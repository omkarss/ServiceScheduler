package customer

import (
	"context"
	"fmt"
	"time"
)

type StandardCustomer struct {
	Customer Customer
}

func NewStandardCustomer(ctx context.Context, fullName string, phoneNumber string, customerType CustomerType, ticketNumber int) (*StandardCustomer, error) {

	return &StandardCustomer{
		Customer{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
			Metadata: &CustomerMetadata{
				TicketNumber: ticketNumber,
				Type:         customerType,
				EntryTime:    time.Now().UTC(),
			},
		},
	}, nil
}

func (*StandardCustomer) GetDetails() (*StandardCustomer, error) {
	return nil, fmt.Errorf("not implemented")
}
