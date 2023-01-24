package server

import (
	"fmt"
	"net/http"

	"github.com/omkar.sunthankar/servicescheduler/internal/deserializer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/customer"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/serializers"
)

func sendResponse(w http.ResponseWriter, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error Serializeing Response ")

	}
}

/*

Checks in Customer
1. Creates a customer
2. Adds customer to neccessary logical queue
3. Returns the customer by serializing it

*/

func (s *Server) CheckInCustomer(w http.ResponseWriter, req *http.Request) {

	var c *customer.Customer
	var vc *customer.VIPCustomer
	var sc *customer.StandardCustomer

	s.Logger.Info("Deserializing request")

	// Deserialize json and validate the input
	c, err := deserializer.DeserializeCustomer(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request: Invalid Request: %s", err), http.StatusInternalServerError)
		return
	}

	s.Logger.Infof("Checking In customer %s", c.FullName)

	switch c.Metadata.Type {
	case customer.CustomerTypeStandard:
		// Check if customer already present in the queue.
		sq := s.Queue[queue.QueueTypeStandard]
		customerExist := sq.Exists(c.FullName, c.PhoneNumber)
		if customerExist {
			http.Error(w, "Customer already present", http.StatusInternalServerError)
			return
		}

		// Create a customer
		sc, err = customer.NewStandardCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeStandard, s.SchedulerTypeA.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}

		// Add to Queue
		sq.Add(&sc.Customer)

		serializedCustomer, err := serializers.SerializeCustomer(&sc.Customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}
		sendResponse(w, serializedCustomer)

	case customer.CustomerTypeVIP:

		vq := s.Queue[queue.QueueTypePriority]
		customerExist := vq.Exists(c.FullName, c.PhoneNumber)
		if customerExist {
			http.Error(w, "Customer already present", http.StatusInternalServerError)
			return
		}

		// Create a customer
		vc, err = customer.NewVIPCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeVIP, s.SchedulerTypeA.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}

		// Add to Queue
		vq.Add(&vc.Customer)

		// Serialize customer to json
		serializedCustomer, err := serializers.SerializeCustomer(&vc.Customer)
		if err != nil {
			http.Error(w, "Error Serializing Response", http.StatusInternalServerError)
			return
		}
		sendResponse(w, serializedCustomer)

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	s.Logger.Infof("Checked in the customer %s ", c.FullName)
}

/*
Scheduler A is a Type A scheduler checks if customer is VIP if yes attends that customer else attends standard customer
1. Uses mutex to make sure we dont attend two customers at the same time
2. Returns  next customer
*/

func (s *Server) GetNextCustomerSchedulerA(w http.ResponseWriter, req *http.Request) {

	s.Logger.Infof("Attending the customer ")

	c, err := s.SchedulerTypeA.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Logger.Infof("Done attending customer %s", c.FullName)

	serializedCustomer, err := serializers.SerializeCustomer(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Serializing response", http.StatusInternalServerError)
		return
	}
	sendResponse(w, serializedCustomer)

}

/*
Scheduler is a Type B scheduler which schedules VIP:Standard customer in ratio 2:1.
In the case when VIP customers are not present and only standard customers are present they schedule standard customer
*/

func (s *Server) GetNextCustomerSchedulerB(w http.ResponseWriter, req *http.Request) {

	s.Logger.Infof("Attending the customer for 10 seconds")

	c, err := s.SchedulerTypeB.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Logger.Infof("Done attending customer %s", c.FullName)

	serializedCustomer, err := serializers.SerializeCustomer(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Serializing response", http.StatusInternalServerError)
		return
	}
	sendResponse(w, serializedCustomer)
}

func (s *Server) GetAllVIPCustomers(w http.ResponseWriter, req *http.Request) {

	vipCustomers := make([]customer.Customer, 0)

	serializedC := make([]byte, 0)

	for _, queue := range s.Queue {

		for _, element := range queue.Elements {

			if element.Metadata.Type == customer.CustomerTypeVIP {
				vipCustomers = append(vipCustomers, *element)
			}
		}
	}

	for _, customer := range vipCustomers {
		serializedCustomer, err := serializers.SerializeCustomer(&customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serializedC = append(serializedC, serializedCustomer...)
	}
	if len(vipCustomers) == 0 {
		http.Error(w, "No VIP customers in the queue", http.StatusInternalServerError)
		return
	}

	sendResponse(w, serializedC)
}

func (s *Server) GetAllStandardCustomers(w http.ResponseWriter, req *http.Request) {

	stdCustomers := make([]customer.Customer, 0)

	serializedC := make([]byte, 0)

	for _, queue := range s.Queue {

		for _, element := range queue.Elements {

			if element.Metadata.Type == customer.CustomerTypeStandard {
				stdCustomers = append(stdCustomers, *element)
			}
		}
	}

	for _, customer := range stdCustomers {
		serializedCustomer, err := serializers.SerializeCustomer(&customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		serializedC = append(serializedC, serializedCustomer...)
	}
	if len(stdCustomers) == 0 {
		http.Error(w, "No standard customers in the queue", http.StatusInternalServerError)
		return
	}

	sendResponse(w, serializedC)
}
