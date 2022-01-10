package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type patchUInt64Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint64 `json:"value"`
}

func ScaleDeploymentReplicas(clientSet *kubernetes.Clientset, namespace string, deploymentname string, scale uint64) error {

	payload := []patchUInt64Value{{
		Op:    "replace",
		Path:  "/spec/replicas",
		Value: scale,
	}}

	payloadbytes, _ := json.Marshal(payload)

	_, err := clientSet.AppsV1().Deployments(namespace).Patch(context.Background(), deploymentname, types.JSONPatchType, payloadbytes, metaV1.PatchOptions{})

	return err
}

func main() {

	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Number of replicas to scale the deployment//
	replicas := flag.Uint64("r", 1, "no of replicas")

	// label of the deployment or deployments that need to be scaled//
	label := flag.String("l", "release=sp-app-blue", "release name label in the format foo=bar")

	// namespace of the deployment or deployments that need to be scaled//
	namespace := flag.String("n", "default", "Namespace of the deployment")

	flag.Parse()

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	// Listing all the deployments with the label

	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metaV1.ListOptions{LabelSelector: *label})
	if err != nil {
		panic(err.Error())
	}

	// For each Deployment in the Deployments listed above (k8s.io/api/apps/v1)
	for i, deployment := range deployments.Items {
		fmt.Printf("Deployment %d: %s is scaled to %d replicas\n", i+1, deployment.ObjectMeta.Name, *replicas)
		err = ScaleDeploymentReplicas(clientset, *namespace, deployment.ObjectMeta.Name, *replicas)
	}
	fmt.Printf("Scaling Completed\n")

	if err != nil {
		panic(err.Error())
	}

}
