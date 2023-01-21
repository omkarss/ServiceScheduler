package server

import (
	"github.com/gorilla/mux"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
)

type Server struct {
	*mux.Router
	VIPFirstSceduler *scheduler.VIPFirstSceduler
	CustomScheduler  *scheduler.CustomScheduler
	Queue            map[queue.QueueType]queue.Queue
}

func NewServer(vipScheduler *scheduler.VIPFirstSceduler, customScheduler *scheduler.CustomScheduler, queue map[queue.QueueType]queue.Queue) *Server {
	return &Server{
		VIPFirstSceduler: vipScheduler,
		CustomScheduler:  customScheduler,
		Queue:            queue,
	}
}

func (s *Server) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/checkIn", s.GetNextCustomer).Methods("POST")
	router.HandleFunc("/getNextCustomer", s.CheckInCustomer).Methods("GET")
}
