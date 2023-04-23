package main

import (
	"context"
	"flag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	dnyClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return
	}
	// Get all pods
	grv := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	_, err = dnyClient.Resource(grv).Namespace("kube-system").Get(context.Background(), "coredns-565d847f94-r27bf", v1.GetOptions{})
	if err != nil {
		return
	}

	deploymentRes := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	//
	//deploymentObject := &unstructured.Unstructured{
	//	Object: map[string]interface{}{
	//		"apiVersion": "apps/v1",
	//		"kind":       "Deployment",
	//		"metadata": map[string]interface{}{
	//			"name": "dest",
	//		},
	//		"spec": map[string]interface{}{
	//			"replicas": 1,
	//			"selector": map[string]interface{}{
	//				"matchLabels": map[string]interface{}{
	//					"app": "test",
	//				},
	//			},
	//			"template": map[string]interface{}{
	//				"metadata": map[string]interface{}{
	//					"labels": map[string]interface{}{
	//						"app": "test",
	//					},
	//				},
	//				"spec": map[string]interface{}{
	//					"containers": []map[string]interface{}{
	//						{
	//							"name":  "test",
	//							"image": "nginx",
	//							"ports": []map[string]interface{}{
	//								{
	//									"name":          "http",
	//									"protocol":      "TCP",
	//									"containerPort": 8080,
	//								},
	//							},
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//output, err := dnyClient.Resource(deploymentRes).Namespace("default").Create(context.Background(), deploymentObject, v1.CreateOptions{})
	//if err != nil {
	//	return
	//}
	//fmt.Print(output)
	//
	dnyClient.Resource(deploymentRes).Namespace("default").Delete(context.Background(), "dest", v1.DeleteOptions{})

}
