package gitlab

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/pager/pkg/apis/devops/v1alpha1"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	git "github.com/xanzy/go-gitlab"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

type userInfo struct {
	username     string
	password     string
	gitlabClient *git.Client
	pagerClient  *pager.Clientset
}

func (u userInfo) Create(obj interface{}) (interface{}, error) {
	user := obj.(*v1alpha2.User)
	users, err := u.list(user.Name)
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		return users[0], nil
	}
	if gitlabUser, _, err := u.gitlabClient.Users.CreateUser(&git.CreateUserOptions{
		Email:               git.String(user.Spec.Email),
		ResetPassword:       git.Bool(true),
		ForceRandomPassword: nil,
		Username:            git.String(user.Name),
		Name:                git.String(user.Name),
		Skype:               nil,
		Linkedin:            nil,
		Twitter:             nil,
		WebsiteURL:          nil,
		Organization:        nil,
		ProjectsLimit:       nil,
		ExternUID:           nil,
		Provider:            nil,
		Bio:                 nil,
		Location:            nil,
		Admin:               nil,
		CanCreateGroup:      git.Bool(false),
		SkipConfirmation:    nil,
		External:            nil,
		PrivateProfile:      nil,
		Note:                nil,
	}); err != nil {
		return nil, err
	} else {
		ctx := context.Background()
		_, err := u.pagerClient.DevopsV1alpha1().Pagers(syncer.DevopsNamespace).Create(ctx, &v1alpha1.Pager{
			ObjectMeta: v1.ObjectMeta{
				Name: "user-" + user.Name,
			},
			Spec: v1alpha1.PagerSpec{
				MessageID:   strconv.Itoa(gitlabUser.ID),
				MessageName: gitlabUser.Name,
				MessageType: user.Kind,
			},
		}, v1.CreateOptions{})
		return gitlabUser, err

	}
}

func (u userInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (u userInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (u userInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (u userInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (u userInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func (u userInfo) list(key string) ([]*git.User, error) {
	users, resp, err := u.gitlabClient.Users.ListUsers(&git.ListUsersOptions{
		Username: git.String(key),
	})
	defer resp.Body.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event":  "list",
			"errros": err.Error(),
			"msg":    resp.Body,
		})
		return nil, err
	} else {
		return users, nil
	}
}

func NewUserGenerator(client *git.Client, pageClient *pager.Clientset) syncer.Generator {
	return &userInfo{
		pagerClient:  pageClient,
		gitlabClient: client,
	}
}
