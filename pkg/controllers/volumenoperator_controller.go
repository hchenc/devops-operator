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
	RegisterReconciler("PersistentVolume", SetUpVolumeReconcile)
}

type VolumeOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (v VolumeOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	volume := &v1.PersistentVolumeClaim{}

	err := v.Get(ctx, req.NamespacedName, volume)
	if err != nil {
		if errors.IsNotFound(err) {
			v.Log.Info("receive delete event")
		}
	} else {
		//sync volume to all environment(fat|uat|smoking)
		_, err := volumeGeneratorService.Add(volume)

		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"event":    "create",
				"resource": "Volume",
				"name":     volume.Name,
				"result":   "failed",
				"error":    err.Error(),
			}).Errorf("volume created failed, retry after %d second", RETRYPERIOD)
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
		log.Logger.WithFields(logrus.Fields{
			"event":    "create",
			"resource": "Volume",
			"name":     volume.Name,
			"result":   "success",
			"message":  "volume controller successful",
		}).Infof("volume %s sync successful", volume.Name)
	}
	return reconcile.Result{}, nil
}

func (v *VolumeOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.PersistentVolumeClaim{}).
		WithEventFilter(&volumePredicate{}).
		Complete(v)
}

type volumePredicate struct {
}

func (v volumePredicate) Create(e event.CreateEvent) bool {
	namespace := e.Meta.GetNamespace()
	if _, exist := e.Meta.GetLabels()["app.kubernetes.io/name"]; !exist {
		return false
	}
	if strings.Contains(namespace, "smoking") || strings.Contains(namespace, "fat") || strings.Contains(namespace, "uat") {
		return true
	} else {
		return false
	}
}

func (v volumePredicate) Delete(event.DeleteEvent) bool {
	return false
}

func (v volumePredicate) Update(event.UpdateEvent) bool {
	return false
}

func (v volumePredicate) Generic(event.GenericEvent) bool {
	return false
}

func SetUpVolumeReconcile(mgr manager.Manager) {
	if err := (&VolumeOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("PersistentVolume"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create volume controller for", err)
	}
}
