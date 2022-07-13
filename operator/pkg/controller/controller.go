package controller

import (
	"os"
	"os/signal"
	"syscall"
	//"log"
	"fmt"
	//"time"
	"context"

	//corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/watch"
	//"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	hzjcs "github.com/ZhengjunHUO/k8s-custom-controller/pkg/client/clientset/versioned"
	//hzjv1alpha1 "github.com/ZhengjunHUO/k8s-custom-controller/pkg/apis/huozj.io/v1alpha1"
)

type Controller struct {
	// clientset for built-in ressource
	clientset	kubernetes.Interface
	// clientset for custom ressource
	hzjclientset	hzjcs.Interface
	// List and watch the delta of certain resource and trigger the event handler
	// normally is to send the key to the queue
	informer	cache.SharedIndexInformer
	hzjinformer	cache.SharedIndexInformer
	// dedicated to the controller to receive event(key)
	queue		workqueue.RateLimitingInterface
	namespace	string
}

func Start(client kubernetes.Interface, hzjclient hzjcs.Interface, namespace string) {
	ctlr := NewController(client, hzjclient, "job", "fufu", namespace)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer fmt.Println("Receive interrupt signal, stop controller, cleanup ...")

	sched := NewScheduler(ctlr, 3, false)
	sched.Start(ctx)
}

func NewController(client kubernetes.Interface, hzjclient hzjcs.Interface, resourceType, customResourceType, namespace string) *Controller {
	// TODO: Implement informers
	return &Controller{}
}

func (c *Controller) Run(ctx context.Context, workersPerCtlr int) error {
	defer utilruntime.HandleCrash()
        //defer c.queue.ShutDown()

	// TODO: Implement informers up
	fmt.Println("Starting controller ...")
	<-ctx.Done()
	fmt.Println("Stopping controller ...")

	return nil
}
