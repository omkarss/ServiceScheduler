package queue

import (
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

type QueueType string
type StandardCustomers []*customer.StandardCustomer
type VIPCustomers []*customer.VIPCustomer

const (
	QueueTypeStandard QueueType = "Standard"
	QueueTypePriority QueueType = "Priority"
)

type QueueI interface {
	Add(c *customer.Customer) (interface{}, error)
	Pop() interface{}
}
