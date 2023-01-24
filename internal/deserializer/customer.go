package deserializer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
)

func DeserializeCustomer(req *http.Request) (*customer.Customer, error) {
	/*
		{
		    "fullName" : "Omkar",
		    "phonenumber": "3159753059",
		    "type": "Priority"
		}
	*/

	var c *Customer

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	t, err := getType(c.Type)
	if err != nil {
		return nil, err
	}

	// Validate
	v := validator.New()
	er := v.Struct(c)
	if er != nil {
		for _, err := range er.(validator.ValidationErrors) {
			return nil, fmt.Errorf("Error validating %s", err.Field())
		}
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

	return "", fmt.Errorf("invalid customerType")
}
