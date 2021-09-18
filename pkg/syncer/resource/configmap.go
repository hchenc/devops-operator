package resource

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type configmapInfo struct {
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func (c configmapInfo) Create(obj interface{}) (interface{}, error) {
	configmap := obj.(*v1.ConfigMap)
	// TODO
	return configmap, nil
}

func (c configmapInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (c configmapInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (c configmapInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (c configmapInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (c configmapInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewConfigmapGenerator(ctx context.Context, kubeClient *kubernetes.Clientset) syncer.Generator {
	return configmapInfo{
		kubeClient: kubeClient,
		ctx:        ctx,
	}
}


