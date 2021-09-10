package gitlab

import (
	"context"
	"github.com/hchenc/application/pkg/apis/app/v1beta1"
	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/devops-operator/pkg/utils"
	"github.com/hchenc/pager/pkg/apis/devops/v1alpha1"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

const (
	CIConfigPath                = ""
	UseCustomTemplate           = true
	TemplateName                = ""
	TemplateProjectID           = 0
	GroupWithProjectTemplatesID = 0
	MergeRequestsTemplate       = ""
	IssuesTemplate              = ""
)

type projectInfo struct {
	projectName      string
	projectNamespace string
	gitlabVersion    string
	gitlabClient     *models.GitlabClient
	pagerClient      *pager.Clientset
	logger           *logrus.Logger
	ctx              context.Context
}

func (p projectInfo) Create(obj interface{}) (interface{}, error) {
	application := obj.(*v1beta1.Application)
	var pip models.Pipelines
	appType := application.Labels["appType"]
	if appType == "" {
		appType = "java"
	}
	for _, pipe := range p.gitlabClient.Pipelines {
		if appType == pipe.Pipeline {
			pip = pipe
		}
	}

	workspaceName := strings.Split(application.Namespace, "-")[0]

	pagerRecord, _ := p.pagerClient.DevopsV1alpha1().Pagers(syncer.DevopsNamespace).Get(p.ctx, "workspace-"+workspaceName, v1.GetOptions{})

	pagerID, _ := strconv.Atoi(pagerRecord.Spec.MessageID)

	name := git.String(application.Name)
	groupID := git.Int(pagerID)
	description := git.String(application.GetAnnotations()["kubesphere.io/description"])

	project, resp, err := p.gitlabClient.Client.Projects.CreateProject(&git.CreateProjectOptions{
		Name:              name,
		Path:              name,
		NamespaceID:       groupID,
		DefaultBranch:     nil,
		Description:       description,
		IssuesAccessLevel: nil,
		//RepositoryAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		//MergeRequestsAccessLevel:            git.AccessControl(git.PrivateAccessControl),
		//ForkingAccessLevel:                  git.AccessControl(git.PrivateAccessControl),
		//BuildsAccessLevel:                   git.AccessControl(git.PrivateAccessControl),
		//WikiAccessLevel:                     git.AccessControl(git.PrivateAccessControl),
		SnippetsAccessLevel: nil,
		PagesAccessLevel:    nil,
		//OperationsAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		EmailsDisabled:                      nil,
		ResolveOutdatedDiffDiscussions:      nil,
		ContainerExpirationPolicyAttributes: nil,
		ContainerRegistryEnabled:            nil,
		SharedRunnersEnabled:                git.Bool(true),
		//Visibility:                          git.Visibility(git.PrivateVisibility),
		ImportURL:                                 nil,
		PublicBuilds:                              nil,
		AllowMergeOnSkippedPipeline:               nil,
		OnlyAllowMergeIfPipelineSucceeds:          nil,
		OnlyAllowMergeIfAllDiscussionsAreResolved: nil,
		MergeMethod:                               nil,
		RemoveSourceBranchAfterMerge:              git.Bool(false),
		LFSEnabled:                                nil,
		RequestAccessEnabled:                      git.Bool(true),
		TagList:                                   nil,
		PrintingMergeRequestLinkEnabled:           nil,
		BuildGitStrategy:                          nil,
		BuildTimeout:                              nil,
		AutoCancelPendingPipelines:                nil,
		BuildCoverageRegex:                        nil,
		CIConfigPath:                              git.String(pip.Ci),
		CIForwardDeploymentEnabled:                nil,
		AutoDevopsEnabled:                         git.Bool(false),
		AutoDevopsDeployStrategy:                  nil,
		ApprovalsBeforeMerge:                      nil,
		ExternalAuthorizationClassificationLabel:  nil,
		Mirror:                                    nil,
		MirrorTriggerBuilds:                       nil,
		InitializeWithReadme:                      git.Bool(true),
		TemplateName:                              git.String(pip.Template),
		//TemplateProjectID:                        git.Int(TemplateProjectID),
		//UseCustomTemplate:                        git.Bool(UseCustomTemplate),
		//GroupWithProjectTemplatesID: git.Int(GroupWithProjectTemplatesID),
		PackagesEnabled:           nil,
		ServiceDeskEnabled:        nil,
		AutocloseReferencedIssues: nil,
		SuggestionCommitMessage:   nil,
		//IssuesTemplate:              git.String(IssuesTemplate),
		//MergeRequestsTemplate:       git.String(MergeRequestsTemplate),
		IssuesEnabled:        git.Bool(true),
		MergeRequestsEnabled: git.Bool(true),
		JobsEnabled:          nil,
		WikiEnabled:          nil,
		SnippetsEnabled:      nil,
	})
	defer resp.Body.Close()
	if err := models.NewConflict(err); err == nil || errors.IsConflict(err) {
		if project == nil {
			if projects, err := p.list(application.Name); err != nil {
				return nil, err
			} else {
				project = projects[0]
			}
		}
		_, err := p.pagerClient.
			DevopsV1alpha1().
			Pagers(syncer.DevopsNamespace).
			Create(p.ctx, &v1alpha1.Pager{
				ObjectMeta: v1.ObjectMeta{
					Name: "application-" + project.Name,
				},
				Spec: v1alpha1.PagerSpec{
					MessageID:   strconv.Itoa(project.ID),
					MessageName: project.Name,
					MessageType: application.Kind,
				},
			}, v1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			return project, nil
		} else {
			return project, err
		}
	}
	return nil, err
}

//func (p projectInfo) assembleProject(name, description *string, groupID *int) *git.CreateProjectOptions {
//	project := &git.CreateProjectOptions{
//		Name:                                name,
//		Path:                                name,
//		NamespaceID:                         groupID,
//		DefaultBranch:                       nil,
//		Description:                         description,
//		IssuesAccessLevel:                   nil,
//		RepositoryAccessLevel:               git.AccessControl(git.PrivateAccessControl),
//		MergeRequestsAccessLevel:            git.AccessControl(git.PrivateAccessControl),
//		ForkingAccessLevel:                  git.AccessControl(git.PrivateAccessControl),
//		BuildsAccessLevel:                   git.AccessControl(git.PrivateAccessControl),
//		WikiAccessLevel:                     git.AccessControl(git.PrivateAccessControl),
//		SnippetsAccessLevel:                 nil,
//		PagesAccessLevel:                    nil,
//		OperationsAccessLevel:               git.AccessControl(git.PrivateAccessControl),
//		EmailsDisabled:                      nil,
//		ResolveOutdatedDiffDiscussions:      nil,
//		ContainerExpirationPolicyAttributes: nil,
//		ContainerRegistryEnabled:            nil,
//		SharedRunnersEnabled:                git.Bool(true),
//		Visibility:                          git.Visibility(git.PrivateVisibility),
//		ImportURL:                           nil,
//		PublicBuilds:                        nil,
//		AllowMergeOnSkippedPipeline:         nil,
//		OnlyAllowMergeIfPipelineSucceeds:    nil,
//		OnlyAllowMergeIfAllDiscussionsAreResolved: nil,
//		MergeMethod:                              nil,
//		RemoveSourceBranchAfterMerge:             git.Bool(false),
//		LFSEnabled:                               nil,
//		RequestAccessEnabled:                     git.Bool(true),
//		TagList:                                  nil,
//		PrintingMergeRequestLinkEnabled:          nil,
//		BuildGitStrategy:                         nil,
//		BuildTimeout:                             nil,
//		AutoCancelPendingPipelines:               nil,
//		BuildCoverageRegex:                       nil,
//		CIConfigPath:                             git.String(CIConfigPath),
//		CIForwardDeploymentEnabled:               nil,
//		AutoDevopsEnabled:                        git.Bool(false),
//		AutoDevopsDeployStrategy:                 nil,
//		ApprovalsBeforeMerge:                     nil,
//		ExternalAuthorizationClassificationLabel: nil,
//		Mirror:                                   nil,
//		MirrorTriggerBuilds:                      nil,
//		InitializeWithReadme:                     git.Bool(true),
//		TemplateName:                             git.String(TemplateName),
//		TemplateProjectID:                        git.Int(TemplateProjectID),
//		UseCustomTemplate:                        git.Bool(UseCustomTemplate),
//		GroupWithProjectTemplatesID:              git.Int(GroupWithProjectTemplatesID),
//		PackagesEnabled:                          nil,
//		ServiceDeskEnabled:                       nil,
//		AutocloseReferencedIssues:                nil,
//		SuggestionCommitMessage:                  nil,
//		IssuesTemplate:                           git.String(IssuesTemplate),
//		MergeRequestsTemplate:                    git.String(MergeRequestsTemplate),
//		IssuesEnabled:                            git.Bool(true),
//		MergeRequestsEnabled:                     git.Bool(true),
//		JobsEnabled:                              nil,
//		WikiEnabled:                              nil,
//		SnippetsEnabled:                          nil,
//	}
//
//	switch p.gitlabVersion {
//	case pipeline.GITLABEEVERSION:
//
//	}
//
//}

func (p projectInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (p projectInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (p projectInfo) GetByName(key string) (interface{}, error) {
	panic("implement me")
}

func (p projectInfo) GetByID(id int) (interface{}, error) {
	//g.gitClient.Groups.GetGroup()
	panic("implement me")
}

func (p projectInfo) List(key string) (interface{}, error) {
	return p.list(key)
}

func (p projectInfo) list(key string) ([]*git.Project, error) {
	projects, resp, err := p.gitlabClient.Client.Projects.ListProjects(&git.ListProjectsOptions{
		Search: git.String(key),
	})
	defer resp.Body.Close()
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"event":  "list",
			"errros": err.Error(),
			"msg":    resp.Body,
		})
		return nil, err
	} else {
		return projects, nil
	}
}

func NewGitLabProjectGenerator(name, group string, ctx context.Context, gitlabClient *models.GitlabClient, pagerClient *pager.Clientset) syncer.Generator {
	logger := utils.GetLogger(logrus.Fields{
		"component": "gitlab.member",
	})
	return &projectInfo{
		projectName:      name,
		projectNamespace: group,
		pagerClient:      pagerClient,
		gitlabClient:     gitlabClient,
		logger:           logger,
		ctx:              ctx,
	}
}
