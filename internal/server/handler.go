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

This handler checks in customer in the queue

*/

func (s *Server) CheckInCustomer(w http.ResponseWriter, req *http.Request) {

	var c *customer.Customer
	var vc *customer.VIPCustomer
	var sc *customer.StandardCustomer

	// Deserialize json and validate the input
	c, err := deserializer.DeserializeCustomer(req)
	if err != nil {
		http.Error(w, "Bad Request: Invalid Request", http.StatusInternalServerError)
		return
	}

	switch c.Metadata.Type {
	case customer.CustomerTypeStandard:
		// Create a customer
		sc, err = customer.NewStandardCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeStandard, s.SchedulerTypeA.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}

		// Add to Queue
		sq := s.Queue[queue.QueueTypeStandard]
		sq.Add(&sc.Customer)

		serializedCustomer, err := serializers.SerializeCustomer(&sc.Customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}
		sendResponse(w, serializedCustomer)

	case customer.CustomerTypeVIP:

		// Create a customer
		vc, err = customer.NewVIPCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeVIP, s.SchedulerTypeA.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}

		// Add to Queue
		sq := s.Queue[queue.QueueTypePriority]
		sq.Add(&vc.Customer)

		// Serialize customer to json
		serializedCustomer, err := serializers.SerializeCustomer(&vc.Customer)
		if err != nil {

		}
		sendResponse(w, serializedCustomer)

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) GetNextCustomerSchedulerA(w http.ResponseWriter, req *http.Request) {

	c, err := s.SchedulerTypeA.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serializedCustomer, err := serializers.SerializeCustomer(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
		return
	}
	sendResponse(w, serializedCustomer)
}

func (s *Server) GetNextCustomerSchedulerB(w http.ResponseWriter, req *http.Request) {

	c, err := s.SchedulerTypeB.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serializedCustomer, err := serializers.SerializeCustomer(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
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
