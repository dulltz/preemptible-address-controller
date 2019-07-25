package controllers

import (
	"context"
	"net"

	networkv1 "github.com/dulltz/preemptible-address-controller/api/v1"
	"github.com/go-logr/logr"
	"google.golang.org/api/compute/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	cName = "kubernetes"
)

// NodeReconciler reconciles a Node object
type NodeReconciler struct {
	client.Client
	Log             logr.Logger
	GCE             *compute.Service
	ProjectID       string
	Region          string
	Zone            string
	AddressLabelKey string
	AddressName     string
	DynamicDNSKey   client.ObjectKey
}

// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.dulltz.com,resources=dynamicdns/status,verbs=get;update;patch

func (r *NodeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log
	var nodeList corev1.NodeList
	err := r.List(ctx, &nodeList)
	if err != nil {
		log.Error(err, "unable to list nodes")
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("listing nodes", "len(nodes)", len(nodeList.Items))
	if len(nodeList.Items) == 0 {
		log.Info("get an empty node list", "project_id", r.ProjectID, "region", r.Region, "name", r.AddressName)
		return ctrl.Result{}, nil
	}

	var aRecords []networkv1.ARecord
	for _, addr := range getNodeExternalAddresses(&nodeList) {
		aRecords = append(aRecords, networkv1.ARecord{Address: addr, CName: cName})
	}
	dd := &networkv1.DynamicDNS{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.DynamicDNSKey.Namespace,
			Name:      r.DynamicDNSKey.Name,
		},
	}
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, dd, func() error {
		dd.Spec.ARecords = aRecords
		return nil
	})
	if err != nil {
		log.Error(err, "unable to mutate DynamicDNS")
		return ctrl.Result{}, ignoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		Owns(&networkv1.DynamicDNS{}).
		Complete(r)
}

func getNodeExternalAddresses(nodeList *corev1.NodeList) []net.IP {
	var res []net.IP
	for _, node := range nodeList.Items {
		for _, a := range node.Status.Addresses {
			if a.Type == corev1.NodeExternalIP {
				res = append(res, net.ParseIP(a.Address))
			}
		}
	}
	return res
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
