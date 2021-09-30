package controller

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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
	RegisterReconciler("ServiceToApp", SetUpServiceReconcile)
}

type ServiceOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (s *ServiceOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	service := &v1.Service{}

	err := s.Get(ctx, req.NamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			s.Log.Info("it's a delete event")
		}
	} else {
		//sync application to all environment(fat|uat|smoking)
		_, err = serviceGeneratorService.Add(service)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "create",
				"resource": "Service",
				"name":     service.Name,
				"result":   "failed",
				"error":    err.Error(),
			}).Errorf("service created failed, retry after %d second", RETRYPERIOD)
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "Service",
			"name":     service.Name,
			"result":   "success",
		}).Infof("service <%s> sync succeed", service.Name)
	}
	return reconcile.Result{}, nil
}

type servicePredicate struct {
}

func (s servicePredicate) Create(e event.CreateEvent) bool {
	name := e.Meta.GetNamespace()
	labels := e.Meta.GetLabels()
	if strings.Contains(name, "smoking") || strings.Contains(name, "fat") || strings.Contains(name, "uat") {
		return true
	} else if _, ok := labels["app"]; ok{
		return true
	} else {
		return false
	}
}
func (s servicePredicate) Update(e event.UpdateEvent) bool {
	//if pod label no changes or add labels, ignore
	return false
}
func (s servicePredicate) Delete(e event.DeleteEvent) bool {
	return false

}
func (s servicePredicate) Generic(e event.GenericEvent) bool {
	return false
}

func (s *ServiceOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Service{}).
		WithEventFilter(&servicePredicate{}).
		Complete(s)
}

func SetUpServiceReconcile(mgr manager.Manager) {
	if err := (&ServiceOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("ServiceToApp"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create service controller for ", err)
	}
}
