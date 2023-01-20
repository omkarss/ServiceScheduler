package customer

import "fmt"

type VIPCustomer struct {
	Customer Customer
}

func (*VIPCustomer) GetDetails() (*VIPCustomer, error) {
	return nil, fmt.Errorf("not implemented")
}
