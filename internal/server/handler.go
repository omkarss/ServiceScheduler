package server

import (
	"encoding/json"
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

func (s *Server) CheckInCustomer(w http.ResponseWriter, req *http.Request) {

	var c *customer.Customer
	var vc *customer.VIPCustomer
	var sc *customer.StandardCustomer

	// Deserialize json and validate the input
	c, err := deserializer.DeserializeCustomer(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request Request")
	}

	switch c.Metadata.Type {
	case customer.CustomerTypeStandard:
		// Create a customer
		sc, err = customer.NewStandardCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeStandard, s.VIPFirstSceduler.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
		}

		// Add to Queue
		sq := s.Queue[queue.QueueTypeStandard]
		sq.Add(&sc.Customer)

		serializedCustomer, err := serializers.SerializeCustomer(sc.Customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}
		sendResponse(w, serializedCustomer)

	case customer.CustomerTypeVIP:

		// Create a customer
		vc, err = customer.NewVIPCustomer(s.ctx, c.FullName, c.PhoneNumber, customer.CustomerTypeVIP, s.VIPFirstSceduler.GetNextTicketNumber(s.ctx))
		if err != nil {
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
		}

		// Add to Queue
		sq := s.Queue[queue.QueueTypePriority]
		sq.Add(&vc.Customer)

		// Serialize customer to json
		serializedCustomer, err := serializers.SerializeCustomer(vc.Customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error Checking In a Customer", http.StatusInternalServerError)
			return
		}
		sendResponse(w, serializedCustomer)

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) GetNextCustomerSchedulerA(w http.ResponseWriter, req *http.Request) {

	c, err := s.VIPFirstSceduler.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
	}
	sendResponse(w, res)
}

func (s *Server) GetNextCustomerSchedulerB(w http.ResponseWriter, req *http.Request) {

	c, err := s.CustomScheduler.GetNextCustomer(s.ctx, s.Queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "Error ", http.StatusInternalServerError)
	}
	sendResponse(w, res)
}

func (s *Server) GetAllCustomers(w http.ResponseWriter, req *http.Request) {

}
