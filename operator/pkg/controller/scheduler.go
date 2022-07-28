package controller

import (
	"context"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var (
	leaseLockName = "hzjoperator"
	leaseDuration = 10*time.Second
	renewDeadline = 5*time.Second
	retryPeriod   = 2*time.Second
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
		s.RunCluster(ctx)
	}else{
		log.Println("Run controller ...")
		s.Run(ctx)
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.ctlr.Run(ctx, s.workersPerCtlr)
}

func (s *Scheduler) RunCluster(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Failed to retrieve node's hostname !")
	}

	llock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: s.ctlr.namespace,
		},
		Client: s.ctlr.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			// unique string identifying a lease holder across all participants in an election
			Identity: hostname,
		},
	}

	lcallbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			log.Println("Starts leading ...")
			s.Run(ctx)
		},
		OnStoppedLeading: func() {
			log.Println("Stops leading.")
		},
		// when the client observes a leader that is not the previously observed leader
		OnNewLeader: func(identity string) {
			if identity == hostname {
				log.Println("Leadership acquired !")
			}else{
				log.Printf("New leader [%s] elected !\n", identity)
			}
		},
	}

	electConf := leaderelection.LeaderElectionConfig{
		Lock:		llock,
		// time to wait before attempting to acquire the lease
		LeaseDuration:	leaseDuration,
		// duration that the acting master will retry refreshing leadership before giving up
		RenewDeadline:	renewDeadline,
		// duration the LeaderElector clients should wait between tries of actions.
		RetryPeriod:	retryPeriod,
		// triggered during certain lifecycle events of the LeaderElector
		Callbacks:	lcallbacks,
		// should be set true if the lock should be released when the run context is cancelled
		ReleaseOnCancel: true,
	}

	// blocks until leader election loop is stopped by ctx, or it has stopped holding the leader lease
	leaderelection.RunOrDie(ctx, electConf)
}
