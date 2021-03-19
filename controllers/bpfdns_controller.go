/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"net"

	"github.com/apex/log"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bpfdnsv1alpha1 "github.com/ValentinoUberti/bpf-dns-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

// BpfdnsReconciler reconciles a Bpfdns object
type BpfdnsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type Ipv4InjectList []string

// +kubebuilder:rbac:groups=bpfdns.bpf.dns,resources=bpfdns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bpfdns.bpf.dns,resources=bpfdns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=bpfdns.bpf.dns,resources=bpfdns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bpfdns object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *BpfdnsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("bpfdns", req.NamespacedName)

	// your logic here
	bpfDNSInstance := &bpfdnsv1alpha1.Bpfdns{}
	workerNodes := &corev1.NodeList{}

	Ipsv4ToInject := make([]string, 0)

	workerSelector := map[string]string{
		"node-role.kubernetes.io/worker": "",
	}

	workerLabelSelector := labels.SelectorFromSet(workerSelector)

	err := r.Client.List(context.TODO(), workerNodes, &client.ListOptions{LabelSelector: workerLabelSelector})

	err = r.Get(ctx, req.NamespacedName, bpfDNSInstance)

	if err != nil {

		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Bpfdns resource not found.")

			return ctrl.Result{}, nil
		} else {

			// Error reading the object - requeue the request.
			log.Error("Failed to get Bpfdns")

		}

		return ctrl.Result{}, nil

	}

	log.Info("Bpfdns resource WAS found.")

	for _, singleNode := range workerNodes.Items {

		log.Info("Node name: " + singleNode.GetName())

	}

	for _, v := range bpfDNSInstance.Spec.BlockDns {
		log.Info(v.DnsName)
		iprecords, _ := net.LookupIP(v.DnsName)
		for _, ip := range iprecords {
			if ipv4 := ip.To4(); ipv4 != nil {

				Ipsv4ToInject = append(Ipsv4ToInject, ipv4.To4().String())

			}
		}
	}

	for _, v := range Ipsv4ToInject {
		fmt.Println("IPv4: ", v)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BpfdnsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bpfdnsv1alpha1.Bpfdns{}).
		Complete(r)
}
