package gitlab

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/pager/pkg/apis/devops/v1alpha1"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

type groupInfo struct {
	gitlabClient *models.GitlabClient
	pagerClient  *pager.Clientset
	logger       *logrus.Logger
	ctx          context.Context
	groupName    string
}

func (g groupInfo) Create(obj interface{}) (interface{}, error) {
	workspace := obj.(*v1alpha2.WorkspaceTemplate)
	name := git.String(workspace.Name)
	description := git.String(workspace.GetAnnotations()["kubesphere.io/description"])
	group, resp, err := g.gitlabClient.Client.Groups.CreateGroup(&git.CreateGroupOptions{
		Name:                           name,
		Path:                           name,
		Description:                    description,
		MembershipLock:                 git.Bool(false),
		Visibility:                     git.Visibility(git.PrivateVisibility),
		ShareWithGroupLock:             git.Bool(false),
		RequireTwoFactorAuth:           git.Bool(false),
		TwoFactorGracePeriod:           nil,
		ProjectCreationLevel:           git.ProjectCreationLevel(git.MaintainerProjectCreation),
		AutoDevopsEnabled:              git.Bool(false),
		SubGroupCreationLevel:          git.SubGroupCreationLevel(git.OwnerSubGroupCreationLevelValue),
		EmailsDisabled:                 git.Bool(false),
		MentionsDisabled:               git.Bool(false),
		LFSEnabled:                     nil,
		RequestAccessEnabled:           nil,
		ParentID:                       nil,
		SharedRunnersMinutesLimit:      nil,
		ExtraSharedRunnersMinutesLimit: nil,
	})
	defer resp.Body.Close()

	if err := models.NewConflict(err); err == nil || errors.IsConflict(err) {
		if group == nil {
			if groups, err := g.list(workspace.Name); err != nil {
				return nil, err
			} else {
				group = groups[0]
			}
		}
		_, err := g.pagerClient.
			DevopsV1alpha1().
			Pagers(syncer.DevopsNamespace).
			Create(g.ctx, &v1alpha1.Pager{
				ObjectMeta: v1.ObjectMeta{
					Name: "workspace-" + workspace.Name,
				},
				Spec: v1alpha1.PagerSpec{
					MessageID:   strconv.Itoa(group.ID),
					MessageName: group.Name,
					MessageType: workspace.Kind,
				},
			}, v1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			return group, nil
		} else {
			return group, err
		}
	} else {
		return nil, err
	}
}

func (g groupInfo) Update(objOld interface{}, objNew interface{}) error {
	if objOld == nil {
		//this is an add operation
	}
	panic("implement me")
}

func (g groupInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (g groupInfo) GetByName(key string) (interface{}, error) {
	return g.list(key)
}

func (g groupInfo) GetByID(id int) (interface{}, error) {
	//g.gitClient.Groups.GetGroup()
	panic("implement me")
}

func (g groupInfo) List(key string) (interface{}, error) {
	groups, err := g.list(key)
	return groups, err
}

func (g groupInfo) list(key string) ([]*git.Group, error) {
	groups, resp, err := g.gitlabClient.Client.Groups.ListGroups(&git.ListGroupsOptions{
		Search: git.String(key),
	})
	defer resp.Body.Close()
	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"event":  "list",
			"errros": err.Error(),
			"msg":    resp.Body,
		})
		return nil, err
	} else {
		return groups, nil
	}
}

func NewGroupGenerator(name string, ctx context.Context, gitlabClient *models.GitlabClient, pagerClient *pager.Clientset) syncer.Generator {
	//cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	return &groupInfo{
		groupName:    name,
		gitlabClient: gitlabClient,
		pagerClient:  pagerClient,
		ctx:          ctx,
	}
}
