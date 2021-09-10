package models

import (
	"context"
	"github.com/hchenc/application/pkg/client/clientset/versioned"
	harbor2 "github.com/hchenc/go-harbor"
	versioned2 "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/xanzy/go-gitlab"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GitlabClient struct {
	Client    *gitlab.Client
	Username  string
	Password  string
	Pipelines []Pipelines
}

type ClientSet struct {
	Ctx context.Context

	Kubeclient *kubernetes.Clientset

	AppClient *versioned.Clientset

	PagerClient *versioned2.Clientset

	GitlabClient *GitlabClient

	HarborClient *harbor2.APIClient
}

func newForDevopsConfigOrDie(devopsConfig *Config) *GitlabClient {
	var gitlabClient GitlabClient

	gc, err := gitlab.NewBasicAuthClient(devopsConfig.Devops.Gitlab.User,
		devopsConfig.Devops.Gitlab.Password,
		gitlab.WithBaseURL("http://"+devopsConfig.Devops.Gitlab.Host+":"+devopsConfig.Devops.Gitlab.Port),
		gitlab.WithoutRetries())
	if err != nil {
		panic(err)
	}

	gitlabClient.Client = gc
	gitlabClient.Username = devopsConfig.Devops.Gitlab.User
	gitlabClient.Password = devopsConfig.Devops.Gitlab.Password
	gitlabClient.Pipelines = devopsConfig.Devops.Pipelines

	return &gitlabClient
}

func NewForConfigOrDie(restConfig *rest.Config, devopsConfig *Config) *ClientSet {
	var cs ClientSet

	cs.Ctx = context.Background()

	cs.Kubeclient = kubernetes.NewForConfigOrDie(restConfig)

	cs.AppClient = versioned.NewForConfigOrDie(restConfig)

	cs.PagerClient = versioned2.NewForConfigOrDie(restConfig)

	cs.GitlabClient = newForDevopsConfigOrDie(devopsConfig)

	cs.HarborClient = harbor2.NewAPIClient(harbor2.NewConfigurationWithContext(devopsConfig.Devops.Harbor.Host,
		context.WithValue(context.Background(),
			harbor2.ContextBasicAuth,
			harbor2.BasicAuth{
				UserName: devopsConfig.Devops.Harbor.User,
				Password: devopsConfig.Devops.Harbor.Password,
			},
		),
	))

	return &cs
}
