package controller

import (
	"context"
	"log"
)

type Scheduler struct {
	ctlr		*Controller
	workersPerCtlr	int
	isCluster	bool
}

func NewScheduler(ctlr *Controller, workersPerCtlr int, isCluster bool) *Scheduler {
	return &Scheduler{
		ctlr: ctlr,
		workersPerCtlr: workersPerCtlr,
		isCluster: isCluster,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	if s.isCluster {
		log.Println("Run controller cluster...")
		// TODO
		// w.RunCluster(ctx)
	}else{
		log.Println("Run controller ...")
		s.Run(ctx)
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.ctlr.Run(ctx, s.workersPerCtlr)
}
