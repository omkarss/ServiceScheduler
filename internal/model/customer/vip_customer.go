package customer

import (
	"context"
	"fmt"
	"time"
)

type VIPCustomer struct {
	Customer Customer
}

func NewVIPCustomer(ctx context.Context, fullName string, phoneNumber string, customerType CustomerType, ticketNumber int) (*VIPCustomer, error) {

	return &VIPCustomer{
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

func (*VIPCustomer) GetDetails() (*VIPCustomer, error) {
	return nil, fmt.Errorf("not implemented")
}
