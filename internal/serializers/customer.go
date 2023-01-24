package serializers

import (
	"encoding/json"
	"fmt"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

func SerializeCustomer(customer *customer.Customer) ([]byte, error) {
	/*
		{
		    "FullName" : "Omkar",
		    "PhoneNumber": "xxxxxx",
		    "Type": "Priority"
		}
	*/

	c := &Customer{
		FullName:     customer.FullName,
		PhoneNumber:  customer.PhoneNumber,
		Type:         string(customer.Metadata.Type),
		TicketNumber: customer.Metadata.TicketNumber,
	}

	jsonC, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("error serializing customer into json ")
	}

	return jsonC, nil
}
