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
	RegisterReconciler("UserToUser", &UserOperatorReconciler{})
}

type UserOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (u *UserOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	user := &iamv1alpha2.User{}

	err := u.Get(ctx, req.NamespacedName, user)
	if err != nil {
		if errors.IsNotFound(err) {
			u.Log.Info("it's a delete event")
		}
	} else {
		// create gitlab project
		_, err := userGeneratorService.Add(user)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
	}
	return reconcile.Result{}, nil
}

func (u *UserOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha2.User{}).
		Complete(u)
}

