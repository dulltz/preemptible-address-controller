package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkv1 "github.com/dulltz/preemptible-address-controller/api/v1"
)

// DynamicDNSReconciler reconciles a DynamicDNS object
type DynamicDNSReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns/status,verbs=get;update;patch

func (r *DynamicDNSReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("dynamicdns", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *DynamicDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.DynamicDNS{}).
		Complete(r)
}
