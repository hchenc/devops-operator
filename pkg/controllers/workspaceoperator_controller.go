package controller

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"

	"time"
)

func init() {
	RegisterReconciler("WorkspaceToGroup", SetUpGroupReconcile)
}

type WorkspaceOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (g *WorkspaceOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	workspaceTemplate := &v1alpha2.WorkspaceTemplate{}

	err := g.Get(ctx, req.NamespacedName, workspaceTemplate)
	if err != nil {
		if errors.IsNotFound(err) {
			g.Log.Info("it's a delete event")
		}
	} else {
		// create gitlab group
		gitlabGroup, err := groupGeneratorService.Add(workspaceTemplate)
		if err != nil {
			if gitlabGroup != nil {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Pager",
					"name":     "workspace-" + workspaceTemplate.Name,
					"result":   "failed",
					"error":    err.Error(),
				}).Errorf("pager created failed, retry after %d second", RETRYPERIOD)
			} else {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Group",
					"name":     workspaceTemplate.Name,
					"result":   "failed",
					"error":    err.Error(),
				}).Errorf("group created failed, retry after %d second", RETRYPERIOD)
			}
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		// create KubeSphere's project(namespace) as environment
		_, err = namespaceGeneratorService.Add(workspaceTemplate)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "create",
				"resource": "Namespace",
				"name":     workspaceTemplate.Name,
				"result":   "failed",
				"error":    err.Error(),
			}).Errorf("namespace created failed, retry after %d second", RETRYPERIOD)
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		// create harbor's project
		_, err = harborGeneratorService.Add(workspaceTemplate)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "create",
				"resource": "Harbor",
				"name":     workspaceTemplate.Name,
				"result":   "failed",
				"error":    err.Error(),
			}).Errorf("harbor project created failed, retry after %d second", RETRYPERIOD)
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "Workspace",
			"name":     workspaceTemplate.Name,
			"result":   "success",
		}).Infof("workspace <%s> sync succeed", workspaceTemplate.Name)
	}
	return reconcile.Result{}, nil
}

func (g *WorkspaceOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.WorkspaceTemplate{}).
		WithEventFilter(&workspacePredicate{}).
		Complete(g)
}

type workspacePredicate struct {
}

func (r workspacePredicate) Create(e event.CreateEvent) bool {
	name := e.Meta.GetName()
	if strings.Contains(name, "system") || strings.Contains(name, "kube") {
		return false
	} else {
		return true
	}
}
func (r workspacePredicate) Update(e event.UpdateEvent) bool {
	//if pod label no changes or add labels, ignore
	return false
}
func (r workspacePredicate) Delete(e event.DeleteEvent) bool {
	return false

}
func (r workspacePredicate) Generic(e event.GenericEvent) bool {
	return false
}

func SetUpGroupReconcile(mgr manager.Manager) {
	if err := (&WorkspaceOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("WorkspaceTemplate"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create workspace controller for", err)
	}
}
