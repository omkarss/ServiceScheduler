package deserializer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

func DeserializeCustomer(req *http.Request) (*customer.Customer, error) {
	/*
		{
		    "FullName" : "Omkar",
		    "PhoneNumber": "3159753059",
		    "Type": "Priority"
		}
	*/

	var c *Customer

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &c)

	t, err := getType(c.Type)
	if err != nil {
		return nil, err
	}

	c1 := &customer.Customer{
		FullName:    c.FullName,
		PhoneNumber: c.PhoneNumber,
		Metadata: &customer.CustomerMetadata{
			Type: t,
		},
	}

	return c1, nil
}

func getType(s string) (customer.CustomerType, error) {
	switch s {
	case "Standard":
		return customer.CustomerTypeStandard, nil

	case "Priority":
		return customer.CustomerTypeVIP, nil

	}

	return "", nil
}
