package controller

import (
	"log"
	"k8s.io/api/core/v1"
	"context"
	hzjv1alpha1 "github.com/ZhengjunHUO/k8s-custom-controller/pkg/apis/huozj.io/v1alpha1"
)

type Handler interface {
	Created(ctx context.Context, item interface{})
	Updated(ctx context.Context, item interface{})
	Deleted(ctx context.Context, item interface{})
}

// default handler simply print out the pod's basic information
type DefaultHandler struct {}

func (dh *DefaultHandler) Created(ctx context.Context, item interface{}) {
	p := item.(*v1.Pod)
	log.Printf("Pod %v/%v created!\n", p.Namespace, p.Name)
}

func (dh *DefaultHandler) Updated(ctx context.Context, item interface{}) {
	p := item.(*v1.Pod)
	log.Printf("Pod %v/%v updated!\n", p.Namespace, p.Name)
}

func (dh *DefaultHandler) Deleted(ctx context.Context, item interface{}) {
	log.Printf("Pod deleted!\n")
}

// print out type fufu's information
type HzjHandler struct {}

func (hh *HzjHandler) Created(ctx context.Context, item interface{}) {
	f := item.(*hzjv1alpha1.Fufu)
	log.Printf("Fufu %v/%v created!\n", f.Namespace, f.Name)
}

func (hh *HzjHandler) Updated(ctx context.Context, item interface{}) {
	f := item.(*hzjv1alpha1.Fufu)
	log.Printf("Fufu %v/%v updated!\n", f.Namespace, f.Name)
}

func (hh *HzjHandler) Deleted(ctx context.Context, item interface{}) {
	log.Printf("Fufu deleted!\n")
}
