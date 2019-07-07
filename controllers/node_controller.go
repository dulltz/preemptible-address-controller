package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/api/compute/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const accessConfigName = "external-nat"

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
}

// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch

func (r *NodeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("req", req.NamespacedName)
	var nodeList corev1.NodeList
	err := r.List(ctx, &nodeList, client.MatchingLabels(map[string]string{r.AddressLabelKey: r.AddressName}))
	if err != nil {
		log.Error(err, "unable to list nodes")
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("listing nodes", "len(nodes)", len(nodeList.Items))
	if len(nodeList.Items) == 0 {
		log.Info("get an empty node list", "project_id", r.ProjectID, "region", r.Region, "name", r.AddressName)
		return ctrl.Result{}, nil
	}

	desiredIP, err := r.getDesiredAddressIP()
	if err != nil {
		log.Error(err, "unable to get address", "project_id", r.ProjectID, "region", r.Region, "name", r.AddressName)
		return ctrl.Result{}, nil
	}

	for _, n := range nodeList.Items {
		for _, a := range n.Status.Addresses {
			if a.Type == corev1.NodeExternalIP && a.Address == desiredIP {
				log.Info("the address is already used", "node", n.Name, "address", desiredIP)
				return ctrl.Result{RequeueAfter: 12 * time.Hour}, nil
			}
		}
	}

	node := nodeList.Items[0]
	log.Info("the address is not used, so add it to the node", "node", node.Name)

	instance, err := r.GCE.Instances.Get(r.ProjectID, r.Zone, node.Name).Do()
	if err != nil {
		log.Error(err, "unable to get the instance", "node", node.Name)
		return ctrl.Result{}, nil
	}

	nicName := r.findNetworkInterfaceName(instance)
	if nicName == "" {
		log.Info(accessConfigName+" not found, so use nic0", "node", node.Name)
		nicName = "nic0"
	}

	ac := &compute.AccessConfig{
		Name:  accessConfigName,
		NatIP: desiredIP,
	}
	_, err = r.GCE.Instances.UpdateAccessConfig(r.ProjectID, r.Zone, node.Name, nicName, ac).Do()
	if err != nil {
		log.Error(err, "unable to update access config to instance", "node", node.Name, "nic", nicName)
		return ctrl.Result{}, nil
	}

	log.Info("update access config to instance", "node", node.Name, "nic", nicName)
	return ctrl.Result{}, nil
}

func (r *NodeReconciler) findNetworkInterfaceName(instance *compute.Instance) string {
	for _, nic := range instance.NetworkInterfaces {
		if len(nic.AccessConfigs) == 0 {
			continue
		}
		if nic.AccessConfigs[0].Name == "external-nat" {
			return nic.Name
		}
	}
	return ""
}

func (r *NodeReconciler) getDesiredAddressIP() (string, error) {
	address, err := r.GCE.Addresses.Get(r.ProjectID, r.Region, r.AddressName).Do()
	if err != nil {
		return "", err
	}
	return address.Address, nil
}

func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		Complete(r)
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
