package server_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

func getQueue(queueType queue.QueueType) *queue.Queue {
	return &queue.Queue{
		Id:       uuid.NewString(),
		Elements: make(queue.CustomerQueue, 0),
		Metadata: &queue.QueueMetadata{
			Type: queueType,
		},
	}
}

func getStandardCustomer(name string, phoneNumber string, cType customer.CustomerType) *customer.StandardCustomer {

	return &customer.StandardCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeStandard,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func getVIPCustomer(name string, phoneNumber string, cType customer.CustomerType) *customer.VIPCustomer {

	return &customer.VIPCustomer{
		Customer: customer.Customer{
			FullName:    name,
			PhoneNumber: phoneNumber,
			Metadata: &customer.CustomerMetadata{
				TicketNumber: 1,
				Type:         customer.CustomerTypeVIP,
				EntryTime:    time.Now().UTC(),
			},
		},
	}
}

func (ts *TestServerSuite) TestCheckInCustomerForVIPFirst() {

	var customer = []byte(`{"FullName":"Omkar","PhoneNumber":"12345","Type":"Standard"}`)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/check-in-customer", bytes.NewBuffer(customer))
	if err != nil {
		ts.T().Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.CheckInCustomer))
	defer dummyServer.Close()
	ts.server.CheckInCustomer(w, req)

	if w.Result().StatusCode != http.StatusOK {
		ts.T().Fatal("status not OK")

	}

	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		ts.T().Fatal("could not serialize")
	}

	// Assertions
	ts.Assert().Equal(customer, body)
	ts.NoError(err)
}

func (ts *TestServerSuite) TestGetNextCustomer() {

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/get-next-customer-schedulerA", nil)
	if err != nil {
		ts.T().Fatal(err)
	}
	dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.GetNextCustomerSchedulerA))
	ts.server.GetNextCustomerSchedulerA(w, req)
	defer dummyServer.Close()

	if w.Result().StatusCode != http.StatusOK {
		ts.T().Fatal("status not OK")
	}

	op, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		ts.T().Fatal("could not serialize")

	}
	fmt.Print(string(op))

}
