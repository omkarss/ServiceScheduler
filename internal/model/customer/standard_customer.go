package customer

import "fmt"

type StandardCustomer struct {
	Customer Customer
}

func (*StandardCustomer) GetDetails() (*StandardCustomer, error) {
	return nil, fmt.Errorf("not implemented")

}
