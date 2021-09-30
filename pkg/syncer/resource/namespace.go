package resource

import (
	"context"
	tenantv1alpha1 "github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha1"
	"github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	FAT     = "功能验收测试环境(Feature Acceptance Test environment)"
	UAT     = "用户验收测试环境(User Acceptance Test environment)"
	SMOKING = "准生产冒烟测试环境(Smoking Test environment)"
)

var (
	env = map[string]string{
		"fat":     FAT,
		"uat":     UAT,
		"smoking": SMOKING,
	}
)

type namespaceInfo struct {
	client *kubernetes.Clientset
	ctx    context.Context
}

func (n namespaceInfo) Create(obj interface{}) (interface{}, error) {
	workspace := obj.(*v1alpha2.WorkspaceTemplate)
	workspaceName := workspace.Name
	candidates := map[string]string{
		"fat":     workspaceName + "-fat",
		"uat":     workspaceName + "-uat",
		"smoking": workspaceName + "-smoking",
	}
	creator := workspace.GetAnnotations()["kubesphere.io/creator"]

	for index, namespaceName := range candidates {
		namespace := assembleResource(workspace, namespaceName, func(obj interface{}, namespace string) interface{} {
			return &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: tenantv1alpha1.SchemeGroupVersion.Group,
							Kind:       "Workspace",
							Name:       workspaceName,
							UID:        workspace.UID,
						},
					},
					Labels: map[string]string{
						"kubesphere.io/creator":       creator,
						"kubernetes.io/metadata.name": namespaceName,
						"kubesphere.io/namespace":     namespaceName,
						"kubesphere.io/workspace":     workspaceName,
					},
					Annotations: map[string]string{
						"kubesphere.io/creator":     creator,
						"kubesphere.io/description": env[index],
					},
				},
			}
		}).(*v1.Namespace)
		_, err := n.client.CoreV1().Namespaces().Create(n.ctx, namespace, metav1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (n namespaceInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (n namespaceInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (n namespaceInfo) GetByName(key string) (interface{}, error) {
	ctx := context.Background()

	ns, err := n.client.CoreV1().Namespaces().Get(ctx, key, metav1.GetOptions{})
	return ns, err
}

func (n namespaceInfo) GetByID(key int) (interface{}, error) {
	panic("implement me")
}

func (n namespaceInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewNamespaceGenerator(ctx context.Context, client *kubernetes.Clientset) syncer.Generator {
	return &namespaceInfo{
		client: client,
		ctx:    ctx,
	}
}
