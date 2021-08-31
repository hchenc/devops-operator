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
	*syncer.ClientSet
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
	for k, v := range candidates{
		if v == application.Namespace{
			continue
		}
		application.Namespace = k
		_, err := a.AppClient.AppV1beta1().Applications(k).Create(ctx, application, v1.CreateOptions{})
		if err != nil{
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

func NewApplicationGenerator(clientset *versioned.Clientset) syncer.Generator {
	return applicationInfo{
		&syncer.ClientSet{
			AppClient: clientset,
		},
	}
}
