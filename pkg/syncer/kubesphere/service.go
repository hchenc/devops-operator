package kubesphere

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type serviceInfo struct {
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func (s serviceInfo) Create(obj interface{}) (interface{}, error) {
	service := obj.(*v1.Service)

	namespacePrefix := strings.Split(service.Namespace, "-")[0]
	candidates := map[string]string{
		namespacePrefix + "-fat":     "fat",
		namespacePrefix + "-uat":     "uat",
		namespacePrefix + "-smoking": "smoking",
	}
	delete(candidates, service.Namespace)

	for namespace := range candidates {
		service := assembleService(service, namespace)
		_, err := s.kubeClient.CoreV1().Services(namespace).Create(s.ctx, service, metav1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (s serviceInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (s serviceInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (s serviceInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (s serviceInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (s serviceInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func assembleService(service *v1.Service, namespace string) *v1.Service {
	return &v1.Service{
		TypeMeta: service.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        service.Name,
			Namespace:   namespace,
			Labels:      service.Labels,
			Annotations: service.Annotations,
			Finalizers:  service.Finalizers,
			ClusterName: service.ClusterName,
		},
		Spec: v1.ServiceSpec{
			Ports:           service.Spec.Ports,
			Selector:        service.Spec.Selector,
			Type:            service.Spec.Type,
			SessionAffinity: service.Spec.SessionAffinity,
		},
	}
}

func NewServiceGenerator(ctx context.Context, kubeClient *kubernetes.Clientset) syncer.Generator {
	return serviceInfo{
		kubeClient: kubeClient,
		ctx:        ctx,
	}
}
