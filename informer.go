package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

func main() {
	// get the kubeconfig path
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//load the kubeconfig from the kubeconfig path
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create k8s client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	stop := make(chan struct{})

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*1)
	fmt.Println("started")
	defer close(stop)
	inf := factory.Core().V1().Pods().Informer()
	inf.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: func(obj interface{}) {

				fmt.Println("Add operation", obj)
			},
			//// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
				newPod := newObj.(*v1.Pod)
				fmt.Println("old", oldPod.ResourceVersion)
				fmt.Println("new", newPod.ResourceVersion)
				if oldPod.ResourceVersion != newPod.ResourceVersion {
					fmt.Println("actucal update")
				}
			},
			//// Called on resource deletion.
			DeleteFunc: func(obj interface{}) {
				fmt.Println("delete operation", obj)
				pod := obj.(*v1.Pod)
				fmt.Println(pod.Name, pod.ObjectMeta.CreationTimestamp)
			},
		})

	go inf.Run(stop)
	<-stop
}

// 1. create informer for configmap,deployments
// 2. update value in configmap and see that change in the informer logs
