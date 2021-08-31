package gitlab

import (
	"context"
	"github.com/hchenc/application/pkg/apis/app/v1beta1"
	"github.com/hchenc/devops-operator/pkg/pipeline"
	"github.com/hchenc/devops-operator/pkg/syncer"
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
	*syncer.ClientSet
}

func (p projectInfo) Create(obj interface{}) (interface{}, error) {
	application := obj.(*v1beta1.Application)

	projects, err := p.list(application.Name)
	if len(projects) != 0 {
		return projects[0], nil
	} else if !errors.IsNotFound(err) {
		return nil, err
	}

	workspaceName := strings.Split(application.Namespace, "")[0]

	ctx := context.Background()
	pagerRecord, _ := p.PagerClient.DevopsV1alpha1().Pagers(syncer.DEVOPS_NAMESPACE).Get(ctx, "workspace-"+workspaceName, v1.GetOptions{})
	pagerID, _ := strconv.Atoi(pagerRecord.Spec.MessageID)

	name := git.String(application.Name)
	groupID := git.Int(pagerID)
	description := git.String(application.GetAnnotations()["kubesphere.io/description"])

	if project, resp, err := p.GitlabClient.Projects.CreateProject(&git.CreateProjectOptions{
		Name:                                name,
		Path:                                name,
		NamespaceID:                         groupID,
		DefaultBranch:                       nil,
		Description:                         description,
		IssuesAccessLevel:                   nil,
		RepositoryAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		MergeRequestsAccessLevel:            git.AccessControl(git.PrivateAccessControl),
		ForkingAccessLevel:                  git.AccessControl(git.PrivateAccessControl),
		BuildsAccessLevel:                   git.AccessControl(git.PrivateAccessControl),
		WikiAccessLevel:                     git.AccessControl(git.PrivateAccessControl),
		SnippetsAccessLevel:                 nil,
		PagesAccessLevel:                    nil,
		OperationsAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		EmailsDisabled:                      nil,
		ResolveOutdatedDiffDiscussions:      nil,
		ContainerExpirationPolicyAttributes: nil,
		ContainerRegistryEnabled:            nil,
		SharedRunnersEnabled:                git.Bool(true),
		Visibility:                          git.Visibility(git.PrivateVisibility),
		ImportURL:                           nil,
		PublicBuilds:                        nil,
		AllowMergeOnSkippedPipeline:         nil,
		OnlyAllowMergeIfPipelineSucceeds:    nil,
		OnlyAllowMergeIfAllDiscussionsAreResolved: nil,
		MergeMethod:                              nil,
		RemoveSourceBranchAfterMerge:             git.Bool(false),
		LFSEnabled:                               nil,
		RequestAccessEnabled:                     git.Bool(true),
		TagList:                                  nil,
		PrintingMergeRequestLinkEnabled:          nil,
		BuildGitStrategy:                         nil,
		BuildTimeout:                             nil,
		AutoCancelPendingPipelines:               nil,
		BuildCoverageRegex:                       nil,
		CIConfigPath:                             git.String(CIConfigPath),
		CIForwardDeploymentEnabled:               nil,
		AutoDevopsEnabled:                        git.Bool(false),
		AutoDevopsDeployStrategy:                 nil,
		ApprovalsBeforeMerge:                     nil,
		ExternalAuthorizationClassificationLabel: nil,
		Mirror:                                   nil,
		MirrorTriggerBuilds:                      nil,
		InitializeWithReadme:                     git.Bool(true),
		TemplateName:                             git.String(TemplateName),
		TemplateProjectID:                        git.Int(TemplateProjectID),
		UseCustomTemplate:                        git.Bool(UseCustomTemplate),
		GroupWithProjectTemplatesID:              git.Int(GroupWithProjectTemplatesID),
		PackagesEnabled:                          nil,
		ServiceDeskEnabled:                       nil,
		AutocloseReferencedIssues:                nil,
		SuggestionCommitMessage:                  nil,
		IssuesTemplate:                           git.String(IssuesTemplate),
		MergeRequestsTemplate:                    git.String(MergeRequestsTemplate),
		IssuesEnabled:                            git.Bool(true),
		MergeRequestsEnabled:                     git.Bool(true),
		JobsEnabled:                              nil,
		WikiEnabled:                              nil,
		SnippetsEnabled:                          nil,
	}); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		ctx := context.Background()
		_, err := p.PagerClient.DevopsV1alpha1().Pagers(syncer.DEVOPS_NAMESPACE).Create(ctx, &v1alpha1.Pager{
			ObjectMeta: v1.ObjectMeta{
				Name: "application-" + project.Name,
			},
			Spec: v1alpha1.PagerSpec{
				MessageID:   strconv.Itoa(project.ID),
				MessageName: project.Name,
				MessageType: application.Kind,
			},
		}, v1.CreateOptions{})
		return project, err
	}
}

func (p projectInfo) assembleProject(name, description *string, groupID *int) *git.CreateProjectOptions{
	project := &git.CreateProjectOptions{
		Name:                                name,
		Path:                                name,
		NamespaceID:                         groupID,
		DefaultBranch:                       nil,
		Description:                         description,
		IssuesAccessLevel:                   nil,
		RepositoryAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		MergeRequestsAccessLevel:            git.AccessControl(git.PrivateAccessControl),
		ForkingAccessLevel:                  git.AccessControl(git.PrivateAccessControl),
		BuildsAccessLevel:                   git.AccessControl(git.PrivateAccessControl),
		WikiAccessLevel:                     git.AccessControl(git.PrivateAccessControl),
		SnippetsAccessLevel:                 nil,
		PagesAccessLevel:                    nil,
		OperationsAccessLevel:               git.AccessControl(git.PrivateAccessControl),
		EmailsDisabled:                      nil,
		ResolveOutdatedDiffDiscussions:      nil,
		ContainerExpirationPolicyAttributes: nil,
		ContainerRegistryEnabled:            nil,
		SharedRunnersEnabled:                git.Bool(true),
		Visibility:                          git.Visibility(git.PrivateVisibility),
		ImportURL:                           nil,
		PublicBuilds:                        nil,
		AllowMergeOnSkippedPipeline:         nil,
		OnlyAllowMergeIfPipelineSucceeds:    nil,
		OnlyAllowMergeIfAllDiscussionsAreResolved: nil,
		MergeMethod:                              nil,
		RemoveSourceBranchAfterMerge:             git.Bool(false),
		LFSEnabled:                               nil,
		RequestAccessEnabled:                     git.Bool(true),
		TagList:                                  nil,
		PrintingMergeRequestLinkEnabled:          nil,
		BuildGitStrategy:                         nil,
		BuildTimeout:                             nil,
		AutoCancelPendingPipelines:               nil,
		BuildCoverageRegex:                       nil,
		CIConfigPath:                             git.String(CIConfigPath),
		CIForwardDeploymentEnabled:               nil,
		AutoDevopsEnabled:                        git.Bool(false),
		AutoDevopsDeployStrategy:                 nil,
		ApprovalsBeforeMerge:                     nil,
		ExternalAuthorizationClassificationLabel: nil,
		Mirror:                                   nil,
		MirrorTriggerBuilds:                      nil,
		InitializeWithReadme:                     git.Bool(true),
		TemplateName:                             git.String(TemplateName),
		TemplateProjectID:                        git.Int(TemplateProjectID),
		UseCustomTemplate:                        git.Bool(UseCustomTemplate),
		GroupWithProjectTemplatesID:              git.Int(GroupWithProjectTemplatesID),
		PackagesEnabled:                          nil,
		ServiceDeskEnabled:                       nil,
		AutocloseReferencedIssues:                nil,
		SuggestionCommitMessage:                  nil,
		IssuesTemplate:                           git.String(IssuesTemplate),
		MergeRequestsTemplate:                    git.String(MergeRequestsTemplate),
		IssuesEnabled:                            git.Bool(true),
		MergeRequestsEnabled:                     git.Bool(true),
		JobsEnabled:                              nil,
		WikiEnabled:                              nil,
		SnippetsEnabled:                          nil,
	}

	switch p.gitlabVersion {
	case pipeline.GITLABEEVERSION:

	}

}

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
	projects, resp, err := p.GitlabClient.Projects.ListProjects(&git.ListProjectsOptions{
		Search: git.String(key),
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
		return projects, nil
	}
}

func NewProjectGenerator(name, group string, gitlabClient *git.Client, pagerClient *pager.Clientset) syncer.Generator {

	return &projectInfo{
		projectName:      name,
		projectNamespace: group,
		ClientSet: &syncer.ClientSet{
			PagerClient:  pagerClient,
			GitlabClient: gitlabClient,
		},
	}
}
