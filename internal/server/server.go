package server

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
)

type Server struct {
	*mux.Router
	ctx              context.Context
	VIPFirstSceduler *scheduler.VIPFirstSceduler
	CustomScheduler  *scheduler.CustomScheduler
	Queue            map[queue.QueueType]*queue.Queue
}

func NewServer(vipScheduler *scheduler.VIPFirstSceduler, customScheduler *scheduler.CustomScheduler, queue map[queue.QueueType]*queue.Queue) *Server {
	s := &Server{
		Router:           mux.NewRouter(),
		ctx:              context.Background(),
		VIPFirstSceduler: vipScheduler,
		CustomScheduler:  customScheduler,
		Queue:            queue,
	}
	s.RegisterRoutes()
	return s
}

func (s *Server) RegisterRoutes() {
	s.HandleFunc("/check-in-customer", s.CheckInCustomer).Methods("POST")
	s.HandleFunc("/get-next-customer-schedulerA", s.GetNextCustomerSchedulerA).Methods("GET")
	s.HandleFunc("/get-next-customer-schedulerB", s.GetNextCustomerSchedulerB).Methods("GET")
}
