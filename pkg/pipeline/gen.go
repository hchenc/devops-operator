package pipeline


import (
	"errors"
	git "github.com/xanzy/go-gitlab"
)

const (
	JAVA_CIFile = `
variables:
  DOCKER_DRIVER: overlay2
  DOCKER_HOST: tcp://localhost:2375

.dind service: &dind_service
  - alias: docker
    name: docker:18.09-dind
    command:
	  - --insecure-registry=

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

func installGitLabClient(host, port, user, password, token string) (*git.Client,error) {
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

func GenGroup(client *git.Client) error {
	groups, resp, err := client.Groups.ListGroups(&git.ListGroupsOptions{
		Search: git.String("devops"),
	})
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if len(groups) != 0{
		return nil
	}

	if _, _, err := client.Groups.CreateGroup(&git.CreateGroupOptions{
		Name:                           git.String("devops"),
		Path:                           git.String("devops"),
		Description:                    git.String("devops"),
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
	}); err != nil {
		return err
	}
	return nil
}

func GenPipeline(client *git.Client) error{
	_,_,err := client.Commits.CreateCommit("",&git.CreateCommitOptions{
		Branch:        git.String(""),
		CommitMessage: git.String("devops pipeline init"),
		StartBranch:   nil,
		StartSHA:      nil,
		StartProject:  nil,
		Actions:       []*git.CommitActionOptions{
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("java.yaml"),
				PreviousPath:    nil,
				Content:         git.String(JAVA_CIFile),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("python.yaml"),
				PreviousPath:    nil,
				Content:         git.String(JAVA_CIFile),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
			{
				Action:          git.FileAction(git.FileCreate),
				FilePath:        git.String("nodejs.yaml"),
				PreviousPath:    nil,
				Content:         git.String(JAVA_CIFile),
				Encoding:        nil,
				LastCommitID:    nil,
				ExecuteFilemode: nil,
			},
		},
		AuthorEmail:   nil,
		AuthorName:    nil,
		Stats:         nil,
		Force:         nil,
	})
	if err != nil{
		return err
	}
	return nil
}

