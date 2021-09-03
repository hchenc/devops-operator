package controller

import (
	"errors"
	app "github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/devops-operator/pkg/syncer/gitlab"
	"github.com/hchenc/devops-operator/pkg/syncer/kubesphere"
	"github.com/hchenc/devops-operator/pkg/utils"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	log = utils.GetLoggerEntry().WithFields(logrus.Fields{
		"component": "controller",
	})
)

var (
	reconcilersMap = make(map[string]Reconciler)

	projectGenerator     syncer.Generator
	groupGenerator       syncer.Generator
	namespaceGenerator   syncer.Generator
	applicationGenerator syncer.Generator
	userGenerator        syncer.Generator
	rolebindingGenerator syncer.Generator
	memberGenerator      syncer.Generator

	projectGeneratorService     syncer.GenerateService
	groupGeneratorService       syncer.GenerateService
	namespaceGeneratorService   syncer.GenerateService
	applicationGeneratorService syncer.GenerateService
	userGeneratorService        syncer.GenerateService
	rolebindingGeneratorService syncer.GenerateService
	memberGeneratorService      syncer.GenerateService
)

type Reconciler interface {
	SetUp(mgr manager.Manager)
}

type Reconcile func(mgr manager.Manager)

func (r Reconcile) SetUp(mgr manager.Manager) {
	r(mgr)
}

func RegisterReconciler(name string, f Reconcile) {
	reconcilersMap[name] = f
}

type DevopsClientet struct {

	kubeclient *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client

}

func (cc *DevopsClientet) Complete(restConfig *rest.Config) {

	cc.kubeclient = kubernetes.NewForConfigOrDie(restConfig)

	cc.appClient = app.NewForConfigOrDie(restConfig)

	cc.pagerClient = pager.NewForConfigOrDie(restConfig)

}

type Controller struct {
	config *models.Config

	kubeclient *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client

	reconcilers map[string]Reconciler

	mgr manager.Manager
}

func (c *Controller) Reconcile(stopCh <-chan struct{}) {
	if err := c.mgr.Start(stopCh); err != nil {
		os.Exit(1)
	}
}

func New(cc *DevopsClientet, mgr manager.Manager, config *models.Config) (*Controller, error) {
	c := &Controller{
		kubeclient:  cc.kubeclient,
		appClient:   cc.appClient,
		pagerClient: cc.pagerClient,
		mgr:         mgr,
	}
	c.reconcilers = reconcilersMap

	var err error
	url := "http://" + config.Devops.Gitlab.Host + ":" + config.Devops.Gitlab.Port
	if config.Devops.Gitlab.Token != "" {
		c.gitlabClient, err = git.NewClient(config.Devops.Gitlab.Token, git.WithBaseURL(url))
		if err != nil {
			return nil, err
		}
	} else if config.Devops.Gitlab.User != "" && config.Devops.Gitlab.Password != "" {
		c.gitlabClient, err = git.NewBasicAuthClient(config.Devops.Gitlab.User, config.Devops.Gitlab.Password, git.WithBaseURL(url))
	} else {
		return nil, errors.New("gitlab certification not provided")
	}

	runtime.Must(installGenerator(c.config, cc.pagerClient, cc.kubeclient, cc.appClient, cc.gitlabClient))
	runtime.Must(installGeneratorService())

	for _, reconciler := range c.reconcilers {
		reconciler.SetUp(mgr)
	}

	return c, nil
}

//func Install(config *rest.Config) {
//	client := kubernetes.NewForConfigOrDie(config)
//
//	appClient := app.NewForConfigOrDie(config)
//
//	pagerClient := pager.NewForConfigOrDie(config)
//
//	runtime.Must(installGenerator(pagerClient, client, appClient))
//	runtime.Must(installGeneratorService())
//}

func installGenerator(config *models.Config, pagerClient *pager.Clientset, clientset *kubernetes.Clientset, appclientset *app.Clientset, gitlabClient *git.Client) error {
	projectGenerator = gitlab.NewProjectGenerator("", "", config, gitlabClient, pagerClient)
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
