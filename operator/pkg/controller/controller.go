package controller

import (
	"os"
	"os/signal"
	"syscall"
	"log"
	"fmt"
	"time"
	"context"

	//corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sinformer "k8s.io/client-go/informers"
	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/watch"
	//"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	hzjcs "github.com/ZhengjunHUO/k8s-custom-controller/pkg/client/clientset/versioned"
	hzjinformer "github.com/ZhengjunHUO/k8s-custom-controller/pkg/client/informers/externalversions"
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

// Sent to queue by informer if match the condition
type Event struct {
        key          string
        eventType    string
        resourceType string
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
	hzjInformerFactory := hzjinformer.NewSharedInformerFactory(hzjclient, 5*time.Second)
	k8sInformerFactory := k8sinformer.NewSharedInformerFactory(client, 5*time.Second)

	// get SharedIndexInformer
	hzjInformer := hzjInformerFactory.Huozj().V1alpha1().Fufus().Informer()
	jobInformer := k8sInformerFactory.Batch().V1().Jobs().Informer()

	var event Event
	var err error
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// register event handler to the informer
	hzjInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event.eventType, event.resourceType = "create", customResourceType
			event.key, err = cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Printf("Resource %v[type %v] created.\n", event.key, event.resourceType)
				queue.Add(event)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			event.eventType, event.resourceType = "update", customResourceType
			event.key, err = cache.MetaNamespaceKeyFunc(oldObj)
			if err == nil {
				log.Printf("Resource %v[type %v] updated.\n", event.key, event.resourceType)
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			event.eventType, event.resourceType = "delete", customResourceType
			event.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Printf("Resource %v[type %v] deleted.\n", event.key, event.resourceType)
				queue.Add(event)
			}
		},
	})

	return &Controller {
		clientset: client,
		hzjclientset: hzjclient,
		informer: jobInformer,
		hzjinformer: hzjInformer,
		queue: queue,
		namespace: namespace,
	}
}

func (c *Controller) Run(ctx context.Context, workersPerCtlr int) error {
	defer utilruntime.HandleCrash()
        //defer c.queue.ShutDown()

	// TODO: Implement informers up
	log.Println("Starting controller ...")
	<-ctx.Done()
	log.Println("Stopping controller ...")

	return nil
}
