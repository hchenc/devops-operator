package kubesphere

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
		"fat":     namespacePrefix + "-fat",
		"uat":     namespacePrefix + "-uat",
		"smoking": namespacePrefix + "-smoking",
	}
	for _, v := range candidates {
		if v == application.Namespace {
			continue
		}
		newApplication := &applicationv1beta1.Application{
			ObjectMeta: v1.ObjectMeta{
				Name:        application.Name,
				Namespace:   v,
				Labels:      application.Labels,
				Annotations: application.Labels,
				Finalizers:  application.Finalizers,
				ClusterName: application.ClusterName,
			},
			Spec: application.Spec,
		}
		_, err := a.appClient.AppV1beta1().Applications(v).Create(a.ctx, newApplication, v1.CreateOptions{})
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
