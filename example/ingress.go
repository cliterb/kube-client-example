package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	apiv1 "k8s.io/api/core/v1"
	exv1beta "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ing struct {
	ingress v1beta1.IngressInterface
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
	var i ing
	i.ingress = clientset.ExtensionsV1beta1().Ingresses(apiv1.NamespaceDefault)
	//i.create_ingress()
	//i.delete_ingress()
	//i.list_ingress()
	//i.watch_ingress()
	i.update_ingress()
}

func (i *ing)create_ingress()  {
	ingress_yaml := &exv1beta.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: exv1beta.IngressSpec{
			Rules: []exv1beta.IngressRule{
				exv1beta.IngressRule{
					Host: "nginx.k8s.local",
					IngressRuleValue: exv1beta.IngressRuleValue{
						HTTP: &exv1beta.HTTPIngressRuleValue{
							Paths: []exv1beta.HTTPIngressPath{
								exv1beta.HTTPIngressPath{
									Backend: exv1beta.IngressBackend{
										ServiceName: "nginx",
										ServicePort: intstr.IntOrString{
											Type: intstr.Int,
											IntVal: 98,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ingress, err := i.ingress.Create(ingress_yaml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ingress %s is created successful", ingress.Name)
}

func (i *ing) delete_ingress ()  {
	err := i.ingress.Delete("nginx",&metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}else {
		fmt.Printf("the ingress delete successful")
	}
}

func (i *ing) list_ingress()  {
	ingress_list, err := i.ingress.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, i1 := range ingress_list.Items {
		fmt.Println(i1.Name)
	}
}

func (i *ing) watch_ingress ()  {
	watch_ingress, err := i.ingress.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	select {
		case e := <-watch_ingress.ResultChan():
			fmt.Println(e.Type)
	}
}

func (i *ing) update_ingress()  {
	ingress_yaml, err := i.ingress.Get("nginx",metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	ingress_yaml.ObjectMeta.Name ="nginx1"
	ingress, err := i.ingress.Update(ingress_yaml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("the %s update successful",ingress.Name)
}
