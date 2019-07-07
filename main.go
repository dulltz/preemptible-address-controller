package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/dulltz/preemptible-address-controller/controllers"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	corev1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var addressLabelKey string
	var addressLabelVal string
	var region string
	var zone string
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&addressLabelKey, "address-label", "preemptible-address", "Label key of the preemptible instance's external address")
	flag.StringVar(&addressLabelVal, "address-name", "", "Name of the target external address")
	flag.StringVar(&region, "region", "us-central1", "Region of the target external address")
	flag.StringVar(&zone, "zone", "us-central1-a", "Zone of the target external address")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	if len(addressLabelVal) == 0 {
		setupLog.Error(errors.New("address-name is required"), "")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := context.Background()
	gce, err := compute.NewService(ctx)
	if err != nil {
		setupLog.Error(err, "unable to create GCE client")
		os.Exit(1)
	}
	credential, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		setupLog.Error(err, "unable to find default credentials")
		os.Exit(1)
	}
	err = (&controllers.NodeReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("Node"),
		GCE:             gce,
		ProjectID:       credential.ProjectID,
		Region:          region,
		Zone:            zone,
		AddressLabelKey: addressLabelKey,
		AddressName:     addressLabelVal,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Node")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
