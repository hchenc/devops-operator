package controller

import (
	"context"
	"github.com/go-logr/logr"
	iamv1alpha2 "github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
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
	RegisterReconciler("UserToUser", SetUpUserReconcile)
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
		gitlabUser, err := userGeneratorService.Add(user)
		if err != nil {
			if gitlabUser != nil {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Pager",
					"name":     "user-" + user.Name,
					"result":   "failed",
					"error":    err.Error(),
				}).Errorf("pager created failed, retry after %d second", RETRYPERIOD)
			} else {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "User",
					"name":     user.Name,
					"result":   "failed",
					"error":    err.Error(),
				}).Errorf("user created failed, retry after %d second", RETRYPERIOD)
			}
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "User",
			"name":     user.Name,
			"result":   "success",
		}).Infof("user <%s> sync succeed", user.Name)
	}
	return reconcile.Result{}, nil
}

type userPredicate struct {
}

func (r userPredicate) Create(e event.CreateEvent) bool {
	name := e.Meta.GetName()
	if strings.Contains(name, "system") || strings.Contains(name, "admin") {
		return false
	} else {
		return true
	}
}
func (r userPredicate) Update(e event.UpdateEvent) bool {
	//if pod label no changes or add labels, ignore
	return false
}
func (r userPredicate) Delete(e event.DeleteEvent) bool {
	return false

}
func (r userPredicate) Generic(e event.GenericEvent) bool {
	return false
}

func (u *UserOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha2.User{}).
		WithEventFilter(&userPredicate{}).
		Complete(u)
}

func SetUpUserReconcile(mgr manager.Manager) {
	if err := (&UserOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("UserToUser"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create user controller for", err)
	}
}
