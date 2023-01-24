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
	var validReq = []byte(`{"FullName":"Omkar","PhoneNumber":"12345","Type":"Standard"}`)
	var invalidReq = []byte(`{"FullName":"Omkar","PhoneNumber":"12345"}`)

	ts.T().Run("Checks in a customer to queue ", func(t *testing.T) {
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
		ts.Assert().Equal(validReq, body)
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
		ts.Assert().Equal(string(op), "Bad Request: Invalid Request\n")
	})

}

func (ts *TestServerSuite) TestGetNextCustomer() {

	c1 := getStandardCustomer("A", "xxxx", customer.CustomerTypeStandard)
	c2 := getVIPCustomer("B", "xxxx", customer.CustomerTypeVIP)

	ts.T().Run("Gets the next Customer from queue", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedc := []byte(`{"FullName":"B","PhoneNumber":"xxxx","Type":"Priority"}`)

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

		ts.Assert().Equal(string(op), string(expectedc))
		ts.Assert().Equal(w.Result().StatusCode, 200)
	})

}
