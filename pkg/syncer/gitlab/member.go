package gitlab

import (
	"context"
	iamv1alpha2 "github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/pager/pkg/apis/devops/v1alpha1"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	git "github.com/xanzy/go-gitlab"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

type memberInfo struct {
	*syncer.ClientSet
}

func (m memberInfo) Create(obj interface{}) (interface{}, error) {
	rolebinding := obj.(*iamv1alpha2.WorkspaceRoleBinding)
	groupName := rolebinding.Labels["kubesphere.io/workspace"]
	userName := rolebinding.Subjects[0].Name

	ctx := context.Background()

	groupRecord, _ := m.PagerClient.DevopsV1alpha1().Pagers(syncer.DEVOPS_NAMESPACE).Get(ctx, "workspace-"+groupName, v1.GetOptions{})

	userRecord, _ := m.PagerClient.DevopsV1alpha1().Pagers(syncer.DEVOPS_NAMESPACE).Get(ctx, "user-"+userName, v1.GetOptions{})

	uid, _ := strconv.Atoi(userRecord.Spec.MessageID)

	if members, resp, err := m.GitlabClient.GroupMembers.AddGroupMember(groupRecord.Spec.MessageID, &git.AddGroupMemberOptions{
		UserID:      git.Int(uid),
		AccessLevel: git.AccessLevel(git.DeveloperPermissions),
		ExpiresAt:   nil,
	}); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		ctx := context.Background()
		_, err := m.PagerClient.DevopsV1alpha1().Pagers(syncer.DEVOPS_NAMESPACE).Create(ctx, &v1alpha1.Pager{
			ObjectMeta: v1.ObjectMeta{
				Name: "member-" + members.Name,
			},
			Spec: v1alpha1.PagerSpec{
				MessageID:   strconv.Itoa(members.ID),
				MessageName: members.Name,
				MessageType: rolebinding.Kind,
			},
		}, v1.CreateOptions{})
		return members, err
	}
}

func (m memberInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (m memberInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (m memberInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (m memberInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (m memberInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewMemberGenerator(gitlabClient *git.Client, pagerClient *pager.Clientset) syncer.Generator {
	return &memberInfo{
		&syncer.ClientSet{
			GitlabClient: gitlabClient,
			PagerClient:  pagerClient,
		},
	}
}
