package controller

import (
	"context"
	"github.com/go-logr/logr"
	iamv1alpha2 "github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func init() {
	RegisterReconciler("RolebindingToMember", &RolebindingOperatorReconciler{})
}

type RolebindingOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r RolebindingOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	rolebinding := &iamv1alpha2.WorkspaceRoleBinding{}

	err := r.Get(ctx, req.NamespacedName, rolebinding)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("it's a delete event")
		}
	} else {
		// add user to group member
		_, err := memberGeneratorService.Add(rolebinding)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		//sync group's user to all environment(fat|uat|smoking)
		_, err = rolebindingGeneratorService.Add(rolebinding)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
	}
	return reconcile.Result{}, nil
}

func (r *RolebindingOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha2.WorkspaceRoleBinding{}).
		Complete(r)
}