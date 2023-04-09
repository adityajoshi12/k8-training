package main

import (
	"context"
	"flag"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
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
	list, err := clientset.AppsV1().Deployments("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println(list.Kind, list)

	//for _, pod := range list.Items {
	//	fmt.Println(pod.Name, pod.Spec)
	//}

	// create pod

	pod := v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "test-pod",
					Image: "nginx",
				},
			},
		},
		Status: v1.PodStatus{},
	}

	podCreated, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metav1.CreateOptions{})
	if err != nil {
		return
	}

	fmt.Println(podCreated)
}

// 1. create deployment
// name: deploy-1
// replicas=2
// namespace = dev - create via code
// 2. Create secret
// username: admin
// password: admin
// 3. Create kubernetes job
// image=apline
// command - echo "hello world"

// client-go https://github.com/feiskyer/go-examples/blob/master/kubernetes
