package main

import (
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type srv struct {
	service corev1.ServiceInterface
}

func main() {
	kubeconfig := flag.String("kubeconfig","./kubeconfig","upload kubeconfig direct")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	var s srv
	s.service = clientset.CoreV1().Services(apiv1.NamespaceDefault)
	//s.create_service()
	//s.delete_service()
	//s.list_service()
	s.update_service()
	s.watch_service()
}

func (s *srv) create_service ()  {
	service_yaml := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []apiv1.ServicePort{
				{
					Name: "nginx",
					Port: 88,
					TargetPort: intstr.IntOrString{
						Type: intstr.Int,
						IntVal: 80,
					},
					Protocol: apiv1.ProtocolTCP,
				},
			},
		},
	}
	service, err := s.service.Create(service_yaml)
	if err != nil {
		panic(err)
	}else {
		fmt.Printf("%s is created successful", service.Name)
	}
}

func (s *srv) delete_service()  {
	err := s.service.Delete("nginx",&metav1.DeleteOptions{})
	if err !=nil {
		panic(err)
	} else {
		fmt.Printf("delete successful")
	}
}

func (s *srv) list_service ()  {
	servicelist, err := s.service.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, i :=range servicelist.Items {
		fmt.Printf("%s \n", i.Name)
	}
}

func (s *srv) update_service()  {
	service_yaml, err := s.service.Get("nginx",metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	service_yaml.Spec.Ports = []apiv1.ServicePort{
		{
			Port: 98,
			TargetPort: intstr.IntOrString{
				Type: intstr.Int,
				IntVal: 80,
			},
			Protocol: apiv1.ProtocolTCP,
		},
	}

	service, err := s.service.Update(service_yaml)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("the service %s  update successful", service.Name)
	}
}

func (s *srv) watch_service()  {
	watch_service, err := s.service.Watch(metav1.ListOptions{})
	if err !=nil {
		panic(err)
	}
	select {
		case e := <-watch_service.ResultChan():
			fmt.Println(e.Type)
	}
}
