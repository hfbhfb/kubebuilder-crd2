/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"time"

	// "github.com/hfbhfb/kubebuilder-crd2/internal/tools/mytrace"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/hfbhfb/kubebuilder-crd2/api/v1alpha1"
	app1v1alpha1 "github.com/hfbhfb/kubebuilder-crd2/api/v1alpha1"
)

var (
	debug1    bool = false
	finalizer      = "foo-finalizer"
)

// FooReconciler reconciles a Foo object
type FooReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app1.crd2.io,resources=foos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app1.crd2.io,resources=foos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app1.crd2.io,resources=foos/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Foo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.2/pkg/reconcile
func (r *FooReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here

	// mytrace.Trace()
	instance := v1alpha1.Foo{}
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Println("11mm")
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if instance.Status.State == v1alpha1.INIT_STATE {
			instance.Status.State = v1alpha1.PENDING_STATE
			r.Status().Update(ctx, &instance)
		}

		// 增加删除时候的处理
		if controllerutil.AddFinalizer(&instance, finalizer) {
			if err := r.Update(ctx, &instance); err != nil {
				return ctrl.Result{}, err
			}
		}

	} else {
		// 等几秒钟后，再删除
		time.Sleep(time.Second * 9)
		controllerutil.RemoveFinalizer(&instance, finalizer)
		if err := r.Update(ctx, &instance); err != nil {
			return ctrl.Result{}, err
		}
		// log.Info("deal finalizer stream: ", finalizer)
	}

	if debug1 {
		fmt.Println("111121aa")
		log.Info("deal delete stream")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FooReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&app1v1alpha1.Foo{}).
		Complete(r)
}
