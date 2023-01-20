package customer

import (
	"context"
	"fmt"
	"time"
)

type CustomerType int

const (
	CustomerTypeStandard CustomerType = iota
	CustomerTypeVIP
)

type CustomerMetadata struct {
	TicketNumber int
	Type         CustomerType
	EntryTime    time.Time
}

type Customer struct {
	FullName    string
	PhoneNumber string
	Metadata    *CustomerMetadata
}

type CustomerI interface {
	GetDetails() (interface{}, error)
}

func NewCustomer(ctx context.Context, fullName string, phoneNumber string, customerType CustomerType, ticketNumber int) (interface{}, error) {

	switch customerType {
	case CustomerTypeStandard:
		return &StandardCustomer{Customer{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
			Metadata: &CustomerMetadata{
				TicketNumber: ticketNumber,
				Type:         customerType,
				EntryTime:    time.Now().UTC(),
			},
		}}, nil

	case CustomerTypeVIP:
		return &VIPCustomer{Customer{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
			Metadata: &CustomerMetadata{
				TicketNumber: ticketNumber,
				Type:         customerType,
				EntryTime:    time.Now().UTC(),
			},
		}}, nil
	default:
		return fmt.Errorf("invalid customer"), nil
	}
}

func (*Customer) GetDetails(*Customer) (*Customer, error) {
	return nil, fmt.Errorf("not implemented")
}
