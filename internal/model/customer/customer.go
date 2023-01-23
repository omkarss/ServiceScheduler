package customer

import (
	"fmt"
	"time"
)

type CustomerType string

const (
	CustomerTypeStandard CustomerType = "Standard"
	CustomerTypeVIP      CustomerType = "Priority"
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

func (*Customer) GetDetails(*Customer) (*Customer, error) {
	return nil, fmt.Errorf("not implemented")
}
