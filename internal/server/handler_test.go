package server_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
)

func (ts *TestServerSuite) TestCheckInCustomerForVIPFirst() {
	var validReq = []byte(`{"FullName":"A","PhoneNumber":"12345","Type":"Standard"}`)
	var invalidReq = []byte(`{"FullName":"A","PhoneNumbe":"12345", "Type": "Standard"}`)

	ts.T().Run("Checks in a customer to queue ", func(t *testing.T) {
		expectedC := []byte(`{"FullName":"A","PhoneNumber":"12345","Type":"Standard","TicketNumber":1}`)
		w := httptest.NewRecorder()
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.CheckInCustomer))
		defer dummyServer.Close()

		req, err := http.NewRequest("POST", "/check-in-customer", bytes.NewBuffer(validReq))
		if err != nil {
			ts.T().Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")

		ts.server.CheckInCustomer(w, req)

		if w.Result().StatusCode != http.StatusOK {
			ts.T().Fatal("status not OK")

		}

		body, err := io.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}

		// Assertions
		ts.Assert().Equal(string(expectedC), string(body))
	})

	ts.T().Run("Fails for Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.CheckInCustomer))
		defer dummyServer.Close()

		req, err := http.NewRequest("POST", "/check-in-customer", bytes.NewBuffer(invalidReq))
		if err != nil {
			ts.T().Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")

		ts.server.CheckInCustomer(w, req)

		op, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}

		// Assertions
		ts.Assert().Equal(w.Result().StatusCode, 500)
		ts.Assert().Equal(string(op), "Bad Request: Invalid Request: Error validating PhoneNumber\n")
	})

}

func (ts *TestServerSuite) TestGetNextCustomerSchedulerA() {

	c1 := getStandardCustomer("A", "xxxx", customer.CustomerTypeStandard, 1)
	c2 := getVIPCustomer("B", "xxxx", customer.CustomerTypeVIP, 2)

	ts.T().Run("Gets the next Customer from queue", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedc := []byte(`{"FullName":"B","PhoneNumber":"xxxx","Type":"Priority","TicketNumber":2}`)

		ts.server.Queue[queue.QueueTypeStandard].Add(&c1.Customer)
		ts.server.Queue[queue.QueueTypePriority].Add(&c2.Customer)

		req, err := http.NewRequest("GET", "/get-next-customer-schedulerA", nil)
		if err != nil {
			ts.T().Fatal(err)
		}
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.GetNextCustomerSchedulerA))
		ts.server.GetNextCustomerSchedulerA(w, req)
		defer dummyServer.Close()

		op, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}

		ts.Assert().Equal(expectedc, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

	})

	ts.T().Run("Gets the next Customer from queue", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedc := []byte(`{"FullName":"A","PhoneNumber":"xxxx","Type":"Standard","TicketNumber":1}`)

		req, err := http.NewRequest("GET", "/get-next-customer-schedulerA", nil)
		if err != nil {
			ts.T().Fatal(err)
		}
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.GetNextCustomerSchedulerA))
		ts.server.GetNextCustomerSchedulerA(w, req)
		defer dummyServer.Close()

		op, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}

		ts.Assert().Equal(expectedc, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

	})

	ts.T().Run("Returns code when there is no customer", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/get-next-customer-schedulerA", nil)
		if err != nil {
			ts.T().Fatal(err)
		}
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.GetNextCustomerSchedulerA))
		ts.server.GetNextCustomerSchedulerA(w, req)
		defer dummyServer.Close()

		op, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}

		ts.Assert().Equal(string(op), "no customer available to attend\n")
		ts.Assert().Equal(w.Result().StatusCode, 500)
	})
}

func (ts *TestServerSuite) TestGetNextCustomerSchedulerB() {

	// Standard Customer
	c1 := getStandardCustomer("A", "xxxx", customer.CustomerTypeStandard, 1)
	c2 := getStandardCustomer("B", "xxxx", customer.CustomerTypeStandard, 2)

	// VIP Customer
	c3 := getVIPCustomer("C", "xxxx", customer.CustomerTypeVIP, 3)
	c4 := getVIPCustomer("D", "xxxx", customer.CustomerTypeVIP, 4)
	c5 := getVIPCustomer("E", "xxxx", customer.CustomerTypeVIP, 5)
	c6 := getVIPCustomer("F", "xxxx", customer.CustomerTypeVIP, 6)

	ts.T().Run("Gets the next Customer from queue", func(t *testing.T) {
		expectedc1 := []byte(`{"FullName":"A","PhoneNumber":"xxxx","Type":"Standard","TicketNumber":1}`)
		expectedc2 := []byte(`{"FullName":"B","PhoneNumber":"xxxx","Type":"Standard","TicketNumber":2}`)
		expectedc3 := []byte(`{"FullName":"C","PhoneNumber":"xxxx","Type":"Priority","TicketNumber":3}`)
		expectedc4 := []byte(`{"FullName":"D","PhoneNumber":"xxxx","Type":"Priority","TicketNumber":4}`)
		expectedc5 := []byte(`{"FullName":"E","PhoneNumber":"xxxx","Type":"Priority","TicketNumber":5}`)
		expectedc6 := []byte(`{"FullName":"F","PhoneNumber":"xxxx","Type":"Priority","TicketNumber":6}`)

		w := httptest.NewRecorder()

		ts.server.Queue[queue.QueueTypeStandard].Add(&c1.Customer)
		ts.server.Queue[queue.QueueTypeStandard].Add(&c2.Customer)
		ts.server.Queue[queue.QueueTypePriority].Add(&c3.Customer)
		ts.server.Queue[queue.QueueTypePriority].Add(&c4.Customer)
		ts.server.Queue[queue.QueueTypePriority].Add(&c5.Customer)
		ts.server.Queue[queue.QueueTypePriority].Add(&c6.Customer)

		req, err := http.NewRequest("GET", "/get-next-customer-schedulerA", nil)
		if err != nil {
			ts.T().Fatal(err)
		}
		dummyServer := httptest.NewServer(http.HandlerFunc(ts.server.GetNextCustomerSchedulerB))
		defer dummyServer.Close()

		// Gets VIP Customer: C
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(expectedc3, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

		w = httptest.NewRecorder()
		// Gets VIP Customer: D
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err = ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(expectedc4, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

		// Gets VIP Customer: A
		w = httptest.NewRecorder()
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err = ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(expectedc1, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

		// Gets VIP Customer: E
		w = httptest.NewRecorder()
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err = ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(expectedc5, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

		// Gets VIP Customer: F
		w = httptest.NewRecorder()
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err = ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(expectedc6, op)
		ts.Assert().Equal(200, w.Result().StatusCode)

		// Gets VIP Customer: B
		w = httptest.NewRecorder()
		ts.server.GetNextCustomerSchedulerB(w, req)
		op, err = ioutil.ReadAll(w.Result().Body)
		if err != nil {
			ts.T().Fatal("could not serialize")
		}
		ts.Assert().Equal(string(expectedc2), string(op))
		ts.Assert().Equal(200, w.Result().StatusCode)

	})
}
