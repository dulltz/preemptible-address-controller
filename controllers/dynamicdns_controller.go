package controllers

import (
	"context"

	networkv1 "github.com/dulltz/preemptible-address-controller/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DynamicDNSReconciler reconciles a DynamicDNS object
type DynamicDNSReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch

func (r *DynamicDNSReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("dynamicdns", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *DynamicDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.DynamicDNS{}).
		Owns(&corev1.Node{}).
		Complete(r)
}
