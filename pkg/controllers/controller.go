package controller

import (
	"errors"
	app "github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/devops-operator/pkg/syncer/gitlab"
	"github.com/hchenc/devops-operator/pkg/syncer/kubesphere"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"os"
)

type Controller struct {

	// Close this to shut down all reconciler.
	StopEverything <-chan struct{}

	client *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client
}

func SetUp(config *restclient.Config) *Controller {

	client := kubernetes.NewForConfigOrDie(config)

	appClient := app.NewForConfigOrDie(config)

	pagerClient := pager.NewForConfigOrDie(config)


	return &Controller{
		StopEverything: nil,
		client:         client,
		appClient:      appClient,
		pagerClient:    pagerClient,
		gitlabClient:   nil,
	}
}


var (
	host     = os.Getenv("HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	token    = ""

	gitlabClient         *git.Client
	projectGenerator     syncer.Generator
	groupGenerator       syncer.Generator
	namespaceGenerator   syncer.Generator
	applicationGenerator syncer.Generator
	userGenerator        syncer.Generator
	rolebindingGenerator syncer.Generator
	memberGenerator 	 syncer.Generator

	projectGeneratorService     syncer.GenerateService
	groupGeneratorService       syncer.GenerateService
	namespaceGeneratorService   syncer.GenerateService
	applicationGeneratorService syncer.GenerateService
	userGeneratorService        syncer.GenerateService
	rolebindingGeneratorService syncer.GenerateService
	memberGeneratorService 	 	syncer.GenerateService
)

func Install(config *restclient.Config) {
	client := kubernetes.NewForConfigOrDie(config)

	appClient := app.NewForConfigOrDie(config)

	pagerClient := pager.NewForConfigOrDie(config)

	runtime.Must(installGitLabClient(host, port, user, password, token))
	runtime.Must(installGenerator(pagerClient, client, appClient))
	runtime.Must(installGeneratorService())
}

func installGitLabClient(host, port, user, password, token string) error {
	var err error
	url := "http://" + host + ":" + port
	if token != "" {
		gitlabClient, err = git.NewClient(token, git.WithBaseURL(url))
		return err
	} else if user != "" && password != "" {
		gitlabClient, err = git.NewBasicAuthClient(user, password, git.WithBaseURL(url))
		return err
	} else {
		return errors.New("gitlab certification not provided")
	}
}

func installGenerator(pagerClient *pager.Clientset, clientset *kubernetes.Clientset, appclientset *app.Clientset) error {
	projectGenerator = gitlab.NewProjectGenerator("", "", gitlabClient, pagerClient)
	groupGenerator = gitlab.NewGroupGenerator("", gitlabClient)
	userGenerator = gitlab.NewUserGenerator(gitlabClient, pagerClient)
	memberGenerator = gitlab.NewMemberGenerator(gitlabClient, pagerClient)

	namespaceGenerator = kubesphere.NewNamespaceGenerator(clientset)
	applicationGenerator = kubesphere.NewApplicationGenerator(appclientset)
	rolebindingGenerator = kubesphere.NewRolebindingGenerator(pagerClient, clientset)
	return nil
}

func installGeneratorService() error {
	projectGeneratorService = syncer.NewGenerateService(projectGenerator)
	groupGeneratorService = syncer.NewGenerateService(groupGenerator)
	namespaceGeneratorService = syncer.NewGenerateService(namespaceGenerator)
	applicationGeneratorService = syncer.NewGenerateService(applicationGenerator)
	userGeneratorService = syncer.NewGenerateService(userGenerator)
	rolebindingGeneratorService = syncer.NewGenerateService(rolebindingGenerator)
	memberGeneratorService = syncer.NewGenerateService(memberGenerator)
	return nil
}

