package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	iamv1alpha2 "github.com/hchenc/devops-operator/pkg/apis/iam/v1alpha2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func init() {
	RegisterReconciler("RolebindingToMember", SetUpRolebindingReconcile)
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
		member, err := memberGeneratorService.Add(rolebinding)
		if err != nil {
			if member != nil {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Pager",
					"name":     "member-" + rolebinding.Name,
					"result":   "failed",
					"error":    err.Error(),
					"message":  fmt.Sprintf("pager created failed, retry after %d second", RETRYPERIOD),
				})
			} else {
				log.Logger.WithFields(logrus.Fields{
					"event":    "create",
					"resource": "Member",
					"name":     rolebinding.Name,
					"result":   "failed",
					"error":    err.Error(),
					"message":  fmt.Sprintf("member created failed, retry after %d second", RETRYPERIOD),
				})
			}
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		//sync group's user to all environment(fat|uat|smoking)
		_, err = rolebindingGeneratorService.Add(rolebinding)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "sync",
				"resource": "Rolebinding",
				"name":     rolebinding.Name,
				"result":   "failed",
				"error":    err.Error(),
				"message":  fmt.Sprintf("rolebinding sync failed, retry after %d second", RETRYPERIOD),
			})
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "Rolebinding",
			"name":     rolebinding.Name,
			"result":   "success",
			"message":  "rolebinding controller successful",
		})
	}
	return reconcile.Result{}, nil
}

type rolebindingPredicate struct {
}

func (r rolebindingPredicate) Create(e event.CreateEvent) bool {
	name := e.Meta.GetName()
	if strings.Contains(name, "system") || strings.Contains(name, "admin") {
		return false
	} else {
		return true
	}
}
func (r rolebindingPredicate) Update(e event.UpdateEvent) bool {
	//if pod label no changes or add labels, ignore
	return false
}
func (r rolebindingPredicate) Delete(e event.DeleteEvent) bool {
	return false

}
func (r rolebindingPredicate) Generic(e event.GenericEvent) bool {
	return false
}

func (r *RolebindingOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha2.WorkspaceRoleBinding{}).
		WithEventFilter(&rolebindingPredicate{}).
		Complete(r)
}

func SetUpRolebindingReconcile(mgr manager.Manager) {
	if err := (&RolebindingOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("RolebindingToMember"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create rolebinding controller for ", err)
	}
}
