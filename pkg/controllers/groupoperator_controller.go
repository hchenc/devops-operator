package controller

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/hchenc/devops-operator/pkg/apis/tenant/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"
)

const (
	RETRYPERIOD = 15
)

type GroupOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (g *GroupOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	workspaceTemplate := &v1alpha2.WorkspaceTemplate{}

	err := g.Get(ctx, req.NamespacedName, workspaceTemplate)
	if err != nil {
		if errors.IsNotFound(err) {
			g.Log.Info("it's a delete event")
		}
	} else {
		// create gitlab group
		_, err := groupGeneratorService.Add(workspaceTemplate)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		// create KubeSphere's project(namespace) as environment
		_, err = namespaceGeneratorService.Add(workspaceTemplate)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
	}
	return reconcile.Result{}, nil
}

func (g *GroupOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
