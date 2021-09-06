package harbor

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
	"github.com/hchenc/devops-operator/pkg/syncer"
	harbor2 "github.com/hchenc/go-harbor"
	"k8s.io/apimachinery/pkg/api/errors"
	"strconv"
)

type projectInfo struct {
	*syncer.ClientSet
	ctx      context.Context
	username string
	password string
	host     string
}

func (p projectInfo) Create(obj interface{}) (interface{}, error) {
	workspace := obj.(*v1alpha2.WorkspaceTemplate)
	if project, err := p.GetByName(workspace.Name); err == nil || errors.IsNotFound(err) {
		return project, nil
	}
	resp, err := p.HarborClient.ProjectApi.CreateProject(p.ctx, harbor2.ProjectReq{
		ProjectName: workspace.Name,
		Metadata: &harbor2.ProjectMetadata{
			Public: "true",
		},
		StorageLimit: 0,
	}, &harbor2.ProjectApiCreateProjectOpts{})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p projectInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (p projectInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (p projectInfo) GetByName(name string) (interface{}, error) {
	project, resp, err := p.HarborClient.ProjectApi.GetProject(p.ctx, name, &harbor2.ProjectApiGetProjectOpts{})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		return project, nil
	}
}

func (p projectInfo) GetByID(id int) (interface{}, error) {
	project, resp, err := p.HarborClient.ProjectApi.GetProject(p.ctx, strconv.Itoa(id), &harbor2.ProjectApiGetProjectOpts{})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		return project, nil
	}
}

func (p projectInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func NewHarborProjectGenerator(name, group string, projectClient *harbor2.APIClient, ctx context.Context) syncer.Generator {
	return &projectInfo{
		ctx: ctx,
		ClientSet: &syncer.ClientSet{
			HarborClient: projectClient,
		},
	}
}
