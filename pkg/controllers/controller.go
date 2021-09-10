package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"os"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
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
	reconcilerMap = make(map[string]Reconciler)

	gitlabProjectGenerator syncer.Generator
	groupGenerator         syncer.Generator
	namespaceGenerator     syncer.Generator
	applicationGenerator   syncer.Generator
	userGenerator          syncer.Generator
	rolebindingGenerator   syncer.Generator
	memberGenerator        syncer.Generator
	harborProjectGenerator syncer.Generator
	deploymentGenerator    syncer.Generator
	serviceGenerator       syncer.Generator

	projectGeneratorService     syncer.GenerateService
	groupGeneratorService       syncer.GenerateService
	namespaceGeneratorService   syncer.GenerateService
	applicationGeneratorService syncer.GenerateService
	userGeneratorService        syncer.GenerateService
	rolebindingGeneratorService syncer.GenerateService
	memberGeneratorService      syncer.GenerateService
	harborGeneratorService      syncer.GenerateService
	deploymentGeneratorService  syncer.GenerateService
	serviceGeneratorService     syncer.GenerateService
)

type Reconciler interface {
	SetUp(mgr manager.Manager)
}

type Reconcile func(mgr manager.Manager)

func (r Reconcile) SetUp(mgr manager.Manager) {
	r(mgr)
}

func RegisterReconciler(name string, f Reconcile) {
	reconcilerMap[name] = f
}

type Controller struct {
	Clientset *models.ClientSet

	ReconcilerMap map[string]Reconciler

	manager manager.Manager
}

func (c *Controller) Reconcile(stopCh <-chan struct{}) {
	if err := c.manager.Start(stopCh); err != nil {
		panic(err)
		os.Exit(1)
	}
}

func NewControllerOrDie(cs *models.ClientSet, mgr manager.Manager) *Controller {
	c := &Controller{
		Clientset: cs,
		manager:   mgr,
	}
	c.ReconcilerMap = reconcilerMap

	runtime.Must(workspace.AddToScheme(mgr.GetScheme()))
	runtime.Must(application.AddToScheme(mgr.GetScheme()))
	runtime.Must(iamv1alpha2.AddToScheme(mgr.GetScheme()))
	runtime.Must(appsv1.AddToScheme(mgr.GetScheme()))
	runtime.Must(corev1.AddToScheme(mgr.GetScheme()))

	installGenerator(c.Clientset)
	installGeneratorService()

	for _, reconciler := range c.ReconcilerMap {
		reconciler.SetUp(mgr)
	}
	return c
}

func installGenerator(clientset *models.ClientSet) {
	gitlabProjectGenerator = gitlab.NewGitLabProjectGenerator("", "", clientset.Ctx, clientset.GitlabClient, clientset.PagerClient)
	groupGenerator = gitlab.NewGroupGenerator("", clientset.Ctx, clientset.GitlabClient, clientset.PagerClient)
	userGenerator = gitlab.NewUserGenerator(clientset.Ctx, clientset.GitlabClient, clientset.PagerClient)
	memberGenerator = gitlab.NewMemberGenerator(clientset.Ctx, clientset.GitlabClient, clientset.PagerClient)

	namespaceGenerator = kubesphere.NewNamespaceGenerator(clientset.Ctx, clientset.Kubeclient)
	applicationGenerator = kubesphere.NewApplicationGenerator(clientset.Ctx, clientset.Kubeclient, clientset.AppClient)
	rolebindingGenerator = kubesphere.NewRolebindingGenerator(clientset.Ctx, clientset.Kubeclient)
	deploymentGenerator = kubesphere.NewDeploymentGenerator(clientset.Ctx, clientset.Kubeclient)
	serviceGenerator = kubesphere.NewServiceGenerator(clientset.Ctx, clientset.Kubeclient)

	harborProjectGenerator = harbor.NewHarborProjectGenerator("", "", clientset.HarborClient)
}

func installGeneratorService() {
	projectGeneratorService = syncer.NewGenerateService(gitlabProjectGenerator)
	groupGeneratorService = syncer.NewGenerateService(groupGenerator)
	namespaceGeneratorService = syncer.NewGenerateService(namespaceGenerator)
	applicationGeneratorService = syncer.NewGenerateService(applicationGenerator)
	userGeneratorService = syncer.NewGenerateService(userGenerator)
	rolebindingGeneratorService = syncer.NewGenerateService(rolebindingGenerator)
	memberGeneratorService = syncer.NewGenerateService(memberGenerator)
	harborGeneratorService = syncer.NewGenerateService(harborProjectGenerator)
	deploymentGeneratorService = syncer.NewGenerateService(deploymentGenerator)
	serviceGeneratorService = syncer.NewGenerateService(serviceGenerator)
}
