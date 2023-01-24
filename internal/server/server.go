package server

import (
	"context"
	"sync"

	"github.com/gorilla/mux"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/queue"
	"github.com/omkar.sunthankar/servicescheduler/internal/model/scheduler"
	"go.uber.org/zap"
)

type Server struct {
	*mux.Router
	mutex          *sync.Mutex
	ctx            context.Context
	Logger         *zap.SugaredLogger
	SchedulerTypeA *scheduler.SchedulerTypeA
	SchedulerTypeB *scheduler.SchedulerTypeB
	Queue          map[queue.QueueType]*queue.Queue
}

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func NewServer(vipScheduler *scheduler.SchedulerTypeA, SchedulerTypeB *scheduler.SchedulerTypeB, queue map[queue.QueueType]*queue.Queue) *Server {
	s := &Server{
		Router:         mux.NewRouter(),
		mutex:          &sync.Mutex{},
		ctx:            context.Background(),
		SchedulerTypeA: vipScheduler,
		SchedulerTypeB: SchedulerTypeB,
		Queue:          queue,
		Logger:         InitLogger(),
	}
	s.RegisterRoutes()
	return s
}

func (s *Server) RegisterRoutes() {
	s.HandleFunc("/check-in-customer", s.CheckInCustomer).Methods("POST")
	s.HandleFunc("/get-next-customer-schedulerA", s.GetNextCustomerSchedulerA).Methods("GET")
	s.HandleFunc("/get-next-customer-schedulerB", s.GetNextCustomerSchedulerB).Methods("GET")
	s.HandleFunc("/get-all-vip-customers", s.GetAllVIPCustomers).Methods("GET")
	s.HandleFunc("/get-all-standard-customers", s.GetAllStandardCustomers).Methods("GET")
}
