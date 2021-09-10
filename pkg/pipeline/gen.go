package pipeline

import (
	"errors"
	"fmt"
	git "github.com/xanzy/go-gitlab"
)

const (
	JAVA_CIFile = `variables:
  DOCKER_DRIVER: overlay2
  DOCKER_HOST: tcp://localhost:2375

.dind service: &dind_service
  - alias: docker
    name: docker:18.09-dind
    command:
      - --insecure-registry=%s

cache: 
  key: mvn-cache
  paths: 
    - .m2/repository

stages:
  - install

install dependency:
  stage: install
  image: maven:3.6.3-jdk-8
  tags:
    - k8s-runner
  script:
    - mvn clean compile -Dmaven.test.skip=true
    - pwd
`
)

func InstallGitLabClient(host, port, user, password, token string) (*git.Client, error) {
	url := "http://" + host + ":" + port
	if token != "" {
		gitlabClient, err := git.NewClient(token, git.WithBaseURL(url))
		return gitlabClient, err
	} else if user != "" && password != "" {
		gitlabClient, err := git.NewBasicAuthClient(user, password, git.WithBaseURL(url))
		return gitlabClient, err
	} else {
		return nil, errors.New("gitlab certification not provided")
	}
}

func GenGroup(client *git.Client) (int, error) {
	groups, resp, err := client.Groups.ListGroups(&git.ListGroupsOptions{
		Search: git.String("devops"),
	})
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}
	if len(groups) != 0 {
		return groups[0].ID, nil
	}

	if gitlabGroup, _, err := client.Groups.CreateGroup(&git.CreateGroupOptions{
		Name:                           git.String("devops"),
		Path:                           git.String("devops"),
		Description:                    git.String("devops"),
		MembershipLock:                 git.Bool(false),
		Visibility:                     git.Visibility(git.InternalVisibility),
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
	}); err != nil {
		return 0, err
	} else {
		return gitlabGroup.ID, nil
	}
}

func GenProject(id int, client *git.Client) (int, error) {
	projects, _, err := client.Projects.ListProjects(&git.ListProjectsOptions{
		Search: git.String("devops"),
	})
	if err != nil {
		return 0, err
	} else if projects != nil {
		return projects[0].ID, nil
	}
	project, _, err := client.Projects.CreateProject(&git.CreateProjectOptions{
		Name:                                git.String("devops"),
		Path:                                git.String("devops"),
		NamespaceID:                         git.Int(id),
		DefaultBranch:                       nil,
		Description:                         nil,
		IssuesAccessLevel:                   nil,
		RepositoryAccessLevel:               nil,
		MergeRequestsAccessLevel:            nil,
		ForkingAccessLevel:                  nil,
		BuildsAccessLevel:                   nil,
		WikiAccessLevel:                     nil,
		SnippetsAccessLevel:                 nil,
		PagesAccessLevel:                    nil,
		OperationsAccessLevel:               nil,
		EmailsDisabled:                      nil,
		ResolveOutdatedDiffDiscussions:      nil,
		ContainerExpirationPolicyAttributes: nil,
		ContainerRegistryEnabled:            nil,
		SharedRunnersEnabled:                nil,
		Visibility:                          nil,
		ImportURL:                           nil,
		PublicBuilds:                        nil,
		AllowMergeOnSkippedPipeline:         nil,
		OnlyAllowMergeIfPipelineSucceeds:    nil,
		OnlyAllowMergeIfAllDiscussionsAreResolved: nil,
		MergeMethod:                              nil,
		RemoveSourceBranchAfterMerge:             nil,
		LFSEnabled:                               nil,
		RequestAccessEnabled:                     nil,
		TagList:                                  nil,
		PrintingMergeRequestLinkEnabled:          nil,
		BuildGitStrategy:                         nil,
		BuildTimeout:                             nil,
		AutoCancelPendingPipelines:               nil,
		BuildCoverageRegex:                       nil,
		CIConfigPath:                             nil,
		CIForwardDeploymentEnabled:               nil,
		AutoDevopsEnabled:                        nil,
		AutoDevopsDeployStrategy:                 nil,
		ApprovalsBeforeMerge:                     nil,
		ExternalAuthorizationClassificationLabel: nil,
		Mirror:                                   nil,
		MirrorTriggerBuilds:                      nil,
		InitializeWithReadme:                     git.Bool(true),
		TemplateName:                             nil,
		TemplateProjectID:                        nil,
		UseCustomTemplate:                        nil,
		GroupWithProjectTemplatesID:              nil,
		PackagesEnabled:                          nil,
		ServiceDeskEnabled:                       nil,
		AutocloseReferencedIssues:                nil,
		SuggestionCommitMessage:                  nil,
		IssuesTemplate:                           nil,
		MergeRequestsTemplate:                    nil,
		IssuesEnabled:                            nil,
		MergeRequestsEnabled:                     nil,
		JobsEnabled:                              nil,
		WikiEnabled:                              nil,
		SnippetsEnabled:                          nil,
	})
	return project.ID, err
}

func GenPipeline(id int, client *git.Client) error {
	_, resp, err := client.Commits.CreateCommit(id, &git.CreateCommitOptions{
		Branch:        git.String("main"),
		CommitMessage: git.String("devops pipeline init"),
		StartBranch:   nil,
		StartSHA:      nil,
		StartProject:  nil,
		Actions: []*git.CommitActionOptions{
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("java.yaml"),
				PreviousPath:    nil,
				Content:         git.String(fmt.Sprintf(JAVA_CIFile, "harbor.hchenc.com")),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("python.yaml"),
				PreviousPath:    nil,
				Content:         git.String(fmt.Sprintf(JAVA_CIFile, "harbor.hchenc.com")),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("nodejs.yaml"),
				PreviousPath:    nil,
				Content:         git.String(fmt.Sprintf(JAVA_CIFile, "harbor.hchenc.com")),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
		},
		AuthorEmail: nil,
		AuthorName:  nil,
		Stats:       nil,
		Force:       nil,
	})
	if resp.StatusCode == 400 {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
