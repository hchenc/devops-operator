package kubesphere

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/syncer"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type rolebindingInfo struct {
	*syncer.ClientSet
}

func (r rolebindingInfo) Create(obj interface{}) (interface{}, error) {
	ctx := context.Background()
	rolebinding := obj.(v1alpha2.WorkspaceRoleBinding)
	workspaceName := rolebinding.Labels["kubesphere.io/workspace"]
	userName := rolebinding.Subjects[0].Name
	candidates := map[string]string{
		"fat":     workspaceName + "-fat",
		"uat":     workspaceName + "-uat",
		"smoking": workspaceName + "-smoking",
	}
	for _, namespace := range candidates {
		_, err := r.Client.RbacV1().RoleBindings(namespace).Get(ctx, userName+"-operator", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			rolebindings := &v1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userName + "-operator",
					Namespace: namespace,
					Labels: map[string]string{
						"iam.kubesphere.io/user-ref": userName,
					},
				},
				Subjects: []v1.Subject{
					{
						Kind:     "User",
						APIGroup: "rbac.authorization.k8s.io",
						Name:     userName,
					},
				},
				RoleRef: v1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Role",
					Name:     "operator",
				},
			}
			_, err := r.Client.RbacV1().RoleBindings(namespace).Create(ctx, rolebindings, metav1.CreateOptions{})
			if err != nil {
				return nil, err
			}
		} else if err == nil {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (r rolebindingInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (r rolebindingInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (r rolebindingInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (r rolebindingInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (r rolebindingInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewRolebindingGenerator(pagerClient *pager.Clientset, clientset *kubernetes.Clientset) syncer.Generator {
	return rolebindingInfo{
		&syncer.ClientSet{
			PagerClient: pagerClient,
			Client:      clientset,
		},
	}
}
