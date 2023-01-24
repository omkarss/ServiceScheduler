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
	TicketNumber int          `json:"TicketNumber" validate:"omitempty"`
	Type         CustomerType `json:"Type" validate:"required" binding:"required"`
	EntryTime    time.Time    `json:"EntryTime" validate:"omitempty"`
}

type Customer struct {
	FullName    string            `json:"FullName" validate:"required" binding:"required"`
	PhoneNumber string            `json:"PhoneNumber" validate:"required" binding:"required"`
	Metadata    *CustomerMetadata `json:"Metadata" validate:"required" binding:"required"`
}

type CustomerI interface {
	GetDetails() (interface{}, error)
}

func (*Customer) GetDetails(*Customer) (*Customer, error) {
	return nil, fmt.Errorf("not implemented")
}
