package controller

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"github.com/hchenc/application/pkg/apis/app/v1beta1"
)

func init() {
	RegisterReconciler("AppToProject", SetUpProjectReconcile)
}

type ProjectOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ProjectOperatorReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	application := &v1beta1.Application{}

	err := r.Get(ctx, req.NamespacedName, application)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("it's a delete event")
		}
	} else {
		// create gitlab project
		_, err := projectGeneratorService.Add(application)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}

		//sync application to all environment(fat|uat|smoking)
		_, err = applicationGeneratorService.Add(application)
		if err != nil {
			return reconcile.Result{
				RequeueAfter: RETRYPERIOD * time.Second,
			}, err
		}
	}
	return reconcile.Result{}, nil
}

func (r *ProjectOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Application{}).
		Complete(r)
}

func SetUpProjectReconcile(mgr manager.Manager) {

	_ = v1beta1.AddToScheme(mgr.GetScheme())

	if err := (&ProjectOperatorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("AppToProject"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr);err != nil{
		log.Fatalf("unable to create project controller")
	}
}
