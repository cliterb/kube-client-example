package main

import (
	"flag"
	"fmt"
	appsbetav1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1beta "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	"encoding/json"
)

func main() {
	kubeconfig := flag.String("kubeconfig","./kubeconfig","absolute path to the kubeconfig file")
	config, err :=clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil {
		panic("the kubeconfig maybe have problem")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("that maybe have problem")
	}
	deploymentclient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	go create_deploy(deploymentclient)
	//go delete_deploy(deploymentclient)
	//go list_deploy(deploymentclient)
	//go update_deploy(deploymentclient)
	watch_deploy(deploymentclient)
}

func create_deploy (deploymentclient v1beta.DeploymentInterface)  {
	var r apiv1.ResourceRequirements
	j := `{"limits": {"cpu":"200m", "memory": "1Gi"}, "requests": {"cpu":"100m", "memory": "100m"}}`
	json.Unmarshal([]byte(j), &r)

	deploy := &appsbetav1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: appsbetav1.DeploymentSpec{
			Replicas: int32Ptr2(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: apiv1.PodSpec{
					Containers:[]apiv1.Container{
						{	Name: "nginx",
							Image: "nginx",
							Resources: r,
							Ports:[]apiv1.ContainerPort{
								{
									Name:  "http",
									Protocol:  apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("that is creating deployment....")
	result, err := deploymentclient.Create(deploy)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success create deployment",result.GetObjectMeta().GetName())
}

func int32Ptr2(i int32) *int32 { return &i }

func list_deploy(deploymentclient v1beta.DeploymentInterface)  {
	deploy, _ := deploymentclient.List(metav1.ListOptions{})
	for _, i := range deploy.Items {
		fmt.Printf("%s have %d replices", i.Name, *i.Spec.Replicas)
	}
}

func delete_deploy(deploymentclient v1beta.DeploymentInterface){
	deletepolicy := metav1.DeletePropagationForeground
	err := deploymentclient.Delete("nginx",&metav1.DeleteOptions{PropagationPolicy:&deletepolicy})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("delete successful")
	}
}

func watch_deploy(deploymentclient v1beta.DeploymentInterface)  {
	w, _ := deploymentclient.Watch(metav1.ListOptions{})
	for {
		select {
			case e := <- w.ResultChan():
				fmt.Println(e.Type,e.Object)
		}
	}
}

func update_deploy(deploymentclient v1beta.DeploymentInterface)  {
	result, err := deploymentclient.Get("nginx",metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	result.Spec.Replicas = int32Ptr2(2)
	deploymentclient.Update(result)
}
