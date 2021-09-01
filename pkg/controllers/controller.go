package controller

import (
	"errors"
	app "github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/config/pipeline"
	"github.com/hchenc/devops-operator/pkg/syncer"
	"github.com/hchenc/devops-operator/pkg/syncer/gitlab"
	"github.com/hchenc/devops-operator/pkg/syncer/kubesphere"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	reconcilers map[string]reconcile.Reconciler

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

func RegisterReconciler(name string, r reconcile.Reconciler) {
	reconcilers[name] = r
}

type CompletedConfig struct {
	config *pipeline.Config

	kubeclient *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client

	reconcilers map[string]reconcile.Reconciler

}

func (cc *CompletedConfig) Complete(config *pipeline.Config, restConfig *rest.Config) error{

	cc.kubeclient = kubernetes.NewForConfigOrDie(restConfig)

	cc.appClient = app.NewForConfigOrDie(restConfig)

	cc.pagerClient = pager.NewForConfigOrDie(restConfig)

	cc.reconcilers = reconcilers

	var err error
	url := "http://" + config.Devops.Gitlab.Host + ":" + config.Devops.Gitlab.Port
	if config.Devops.Gitlab.Token != "" {
		cc.gitlabClient, err = git.NewClient(config.Devops.Gitlab.Token, git.WithBaseURL(url))
		if err != nil{
			return err
		}
	} else if config.Devops.Gitlab.User != "" && config.Devops.Gitlab.Password != "" {
		cc.gitlabClient, err = git.NewBasicAuthClient(config.Devops.Gitlab.User, config.Devops.Gitlab.Password, git.WithBaseURL(url))
	} else {
		return errors.New("gitlab certification not provided")
	}
	return err
}

type Controller struct {

	config *pipeline.Config

	kubeclient *kubernetes.Clientset

	appClient *app.Clientset

	pagerClient *pager.Clientset

	gitlabClient *git.Client

	reconcilers map[string]reconcile.Reconciler

	mgr controllerruntime.Manager
}

func (c *Controller) Reconcile(stopCh chan struct{}){
	if err := c.mgr.Start(stopCh); err != nil {
		os.Exit(1)
	}
}

func New(cc *CompletedConfig, mgr controllerruntime.Manager) (*Controller, error) {

	c := &Controller{
		kubeclient:   cc.kubeclient,
		appClient:    cc.appClient,
		pagerClient:  cc.pagerClient,
		gitlabClient: cc.gitlabClient,
		mgr:          mgr,
	}
	c.reconcilers = cc.reconcilers
	for _,reconciler := range reconcilers{
		reconciler
	}

}

//type Runable interface {
//	Run(stopCh chan struct{})
//}
//
//type Reconciler func(stopCh chan struct{})
//
//func (r Reconciler) Run(stopCh chan struct{})  {
//	r(stopCh)
//}
//
//var ReconcilersMap map[string]Reconciler
//
//func RegisterReconciler(name string, f Reconciler) {
//	ReconcilersMap[name] = f
//}




func Install(config *rest.Config) {
	client := kubernetes.NewForConfigOrDie(config)

	appClient := app.NewForConfigOrDie(config)

	pagerClient := pager.NewForConfigOrDie(config)

	runtime.Must(installGenerator(pagerClient, client, appClient))
	runtime.Must(installGeneratorService())
}

//
//func installGitLabClient(host, port, user, password, token string) error {
//	var err error
//	url := "http://" + host + ":" + port
//	if token != "" {
//		gitlabClient, err = git.NewClient(token, git.WithBaseURL(url))
//		return err
//	} else if user != "" && password != "" {
//		gitlabClient, err = git.NewBasicAuthClient(user, password, git.WithBaseURL(url))
//		return err
//	} else {
//		return errors.New("gitlab certification not provided")
//	}
//}

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

