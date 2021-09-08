package kubesphere

import (
	"context"
	applicationv1beta1 "github.com/hchenc/application/pkg/apis/app/v1beta1"
	"github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type applicationInfo struct {
	appClient *versioned.Clientset
}

func (a applicationInfo) Create(obj interface{}) (interface{}, error) {
	ctx := context.Background()
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
		if exist, _ := a.appClient.AppV1beta1().Applications(v).Get(ctx,application.Name, v1.GetOptions{});exist != nil {
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
			Spec:application.Spec,
		}
		_, err := a.appClient.AppV1beta1().Applications(v).Create(ctx, newApplication, v1.CreateOptions{})
		if err != nil {
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

func NewApplicationGenerator(appClient *versioned.Clientset) syncer.Generator {
	return applicationInfo{
		appClient: appClient,
	}
}
