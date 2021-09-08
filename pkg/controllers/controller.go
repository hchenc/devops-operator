package controller

import (
	"context"
	"errors"
	"os"

	app "github.com/hchenc/application/pkg/client/clientset/versioned"
	harbor2 "github.com/hchenc/go-harbor"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/devops-operator/pkg/syncer/gitlab"
	"github.com/hchenc/devops-operator/pkg/syncer/harbor"
	"github.com/hchenc/devops-operator/pkg/syncer/kubesphere"
	"github.com/hchenc/devops-operator/pkg/utils"

	application "github.com/hchenc/application/pkg/apis/app/v1beta1"
	iamv1alpha2 "github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	workspace "github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
)

var (
	log = utils.GetLoggerEntry().WithFields(logrus.Fields{
		"component": "controller",
	})
)

var (
	reconcilersMap = make(map[string]Reconciler)

	gitlabProjectGenerator syncer.Generator
	groupGenerator         syncer.Generator
	namespaceGenerator     syncer.Generator
	applicationGenerator   syncer.Generator
	userGenerator          syncer.Generator
	rolebindingGenerator   syncer.Generator
	memberGenerator        syncer.Generator
	harborProjectGenerator syncer.Generator

	projectGeneratorService     syncer.GenerateService
	groupGeneratorService       syncer.GenerateService
	namespaceGeneratorService   syncer.GenerateService
	applicationGeneratorService syncer.GenerateService
	userGeneratorService        syncer.GenerateService
	rolebindingGeneratorService syncer.GenerateService
	memberGeneratorService      syncer.GenerateService
	harborGeneratorService      syncer.GenerateService
)

type ClientSet struct {
	config *models.Config

	kubeclient *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client

	harborClient *harbor2.APIClient
}

func (cs *ClientSet) Initial(restConfig *rest.Config, devopsConfig *models.Config) {

	cs.config = devopsConfig

	cs.kubeclient = kubernetes.NewForConfigOrDie(restConfig)

	cs.appClient = app.NewForConfigOrDie(restConfig)

	cs.pagerClient = pager.NewForConfigOrDie(restConfig)

	var err error
	url := "http://" + devopsConfig.Devops.Gitlab.Host + ":" + devopsConfig.Devops.Gitlab.Port
	if devopsConfig.Devops.Gitlab.Token != "" {
		cs.gitlabClient, err = git.NewClient(devopsConfig.Devops.Gitlab.Token, git.WithBaseURL(url))
		if err != nil {
			panic(err)
		}
	} else if devopsConfig.Devops.Gitlab.User != "" && devopsConfig.Devops.Gitlab.Password != "" {
		cs.gitlabClient, err = git.NewBasicAuthClient(devopsConfig.Devops.Gitlab.User, devopsConfig.Devops.Gitlab.Password, git.WithBaseURL(url))
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("gitlab certification not provided"))
	}

	harborCfg := harbor2.NewConfigurationWithContext(devopsConfig.Devops.Harbor.Host, context.WithValue(context.Background(), harbor2.ContextBasicAuth, harbor2.BasicAuth{
		UserName: devopsConfig.Devops.Harbor.User,
		Password: devopsConfig.Devops.Harbor.Password,
	}))

	cs.harborClient = harbor2.NewAPIClient(harborCfg)

}

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

type Controller struct {
	clientset *ClientSet

	reconcilers map[string]Reconciler

	mgr manager.Manager
}

func (c *Controller) Reconcile(stopCh <-chan struct{}) {
	if err := c.mgr.Start(stopCh); err != nil {
		os.Exit(1)
	}
}

func New(cs *ClientSet, mgr manager.Manager) (*Controller, error) {
	c := &Controller{
		clientset: cs,
		mgr:       mgr,
	}
	c.reconcilers = reconcilersMap

	runtime.Must(workspace.AddToScheme(mgr.GetScheme()))
	runtime.Must(application.AddToScheme(mgr.GetScheme()))
	runtime.Must(iamv1alpha2.AddToScheme(mgr.GetScheme()))

	runtime.Must(installGenerator(c))
	runtime.Must(installGeneratorService())

	for _, reconciler := range c.reconcilers {
		reconciler.SetUp(mgr)
	}

	return c, nil
}

func installGenerator(c *Controller) error {
	gitlabProjectGenerator = gitlab.NewGitLabProjectGenerator("", "", c.clientset.config, c.clientset.gitlabClient, c.clientset.pagerClient)
	groupGenerator = gitlab.NewGroupGenerator("", c.clientset.gitlabClient, c.clientset.pagerClient)
	userGenerator = gitlab.NewUserGenerator(c.clientset.gitlabClient, c.clientset.pagerClient)
	memberGenerator = gitlab.NewMemberGenerator(c.clientset.gitlabClient, c.clientset.pagerClient)

	namespaceGenerator = kubesphere.NewNamespaceGenerator(c.clientset.kubeclient)
	applicationGenerator = kubesphere.NewApplicationGenerator(c.clientset.appClient)
	rolebindingGenerator = kubesphere.NewRolebindingGenerator(c.clientset.kubeclient)

	harborProjectGenerator = harbor.NewHarborProjectGenerator("", "", c.clientset.harborClient)

	return nil
}

func installGeneratorService() error {
	projectGeneratorService = syncer.NewGenerateService(gitlabProjectGenerator)
	groupGeneratorService = syncer.NewGenerateService(groupGenerator)
	namespaceGeneratorService = syncer.NewGenerateService(namespaceGenerator)
	applicationGeneratorService = syncer.NewGenerateService(applicationGenerator)
	userGeneratorService = syncer.NewGenerateService(userGenerator)
	rolebindingGeneratorService = syncer.NewGenerateService(rolebindingGenerator)
	memberGeneratorService = syncer.NewGenerateService(memberGenerator)
	harborGeneratorService = syncer.NewGenerateService(harborProjectGenerator)
	return nil
}
