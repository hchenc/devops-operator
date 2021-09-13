package resource

import (
	"context"
	applicationv1beta1 "github.com/hchenc/application/pkg/apis/app/v1beta1"
	"github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type applicationInfo struct {
	appClient  *versioned.Clientset
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func (a applicationInfo) Create(obj interface{}) (interface{}, error) {
	application := obj.(*applicationv1beta1.Application)
	namespacePrefix := strings.Split(application.Namespace, "-")[0]
	candidates := map[string]string{
		namespacePrefix + "-fat":     "fat",
		namespacePrefix + "-uat":     "uat",
		namespacePrefix + "-smoking": "smoking",
	}
	delete(candidates, application.Namespace)

	for namespace, _ := range candidates {
		application := assembleResource(application, namespace, func(obj interface{}, namespace string) interface{} {
			return &applicationv1beta1.Application{
				ObjectMeta: v1.ObjectMeta{
					Name:        application.Name,
					Namespace:   namespace,
					Labels:      application.Labels,
					Annotations: application.Labels,
					Finalizers:  application.Finalizers,
					ClusterName: application.ClusterName,
				},
				Spec: application.Spec,
			}
		}).(*applicationv1beta1.Application)
		_, err := a.appClient.AppV1beta1().Applications(namespace).Create(a.ctx, application, v1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (a applicationInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (a applicationInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (a applicationInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (a applicationInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (a applicationInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewApplicationGenerator(ctx context.Context, kubeClient *kubernetes.Clientset, appClient *versioned.Clientset) syncer.Generator {
	return applicationInfo{
		appClient:  appClient,
		kubeClient: kubeClient,
		ctx:        ctx,
	}
}
