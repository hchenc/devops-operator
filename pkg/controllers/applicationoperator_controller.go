package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
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

	"github.com/hchenc/application/pkg/apis/app/v1beta1"
)

func init() {
	RegisterReconciler("AppToProject", SetUpProjectReconcile)
}

type ApplicationOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ApplicationOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	application := &v1beta1.Application{}

	err := r.Get(ctx, req.NamespacedName, application)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("it's a delete event")
		}
	} else {
		// create gitlab project
		project, err := projectGeneratorService.Add(application)
		if err != nil {
			if project != nil {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Pager",
					"name":     "application-" + application.Name,
					"result":   "failed",
					"error":    err.Error(),
					"message":  fmt.Sprintf("pager created failed, retry after %d second", RETRYPERIOD),
				})
			} else {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Project",
					"name":     application.Name,
					"result":   "failed",
					"error":    err.Error(),
					"message":  fmt.Sprintf("project created failed, retry after %d second", RETRYPERIOD),
				})
			}
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		//sync application to all environment(fat|uat|smoking)
		_, err = applicationGeneratorService.Add(application)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "create",
				"resource": "Application",
				"name":     application.Name,
				"result":   "failed",
				"error":    err.Error(),
				"message":  fmt.Sprintf("application created failed, retry after %d second", RETRYPERIOD),
			})
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "Application",
			"name":     application.Name,
			"result":   "success",
			"message":  "application controller successful",
		})
	}
	return reconcile.Result{}, nil
}

type projectPredicate struct {
}

func (r projectPredicate) Create(e event.CreateEvent) bool {
	name := e.Meta.GetName()
	if strings.Contains(name, "system") || strings.Contains(name, "kube") {
		return false
	} else {
		return true
	}
}
func (r projectPredicate) Update(e event.UpdateEvent) bool {
	//if pod label no changes or add labels, ignore
	return false
}
func (r projectPredicate) Delete(e event.DeleteEvent) bool {
	return false

}
func (r projectPredicate) Generic(e event.GenericEvent) bool {
	return false
}

func (r *ApplicationOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Application{}).
		WithEventFilter(&projectPredicate{}).
		Complete(r)
}

func SetUpProjectReconcile(mgr manager.Manager) {
	if err := (&ApplicationOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("AppToProject"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create application controller for ", err)
	}
}
