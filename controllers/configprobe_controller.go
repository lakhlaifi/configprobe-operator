/*
Copyright 2023.

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	syntheticv1 "github.com/lakhlaifi/configprobe-operator/api/v1"
)

// ConfigProbeReconciler reconciles a ConfigProbe object
type ConfigProbeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=synthetic.clodevo.com,resources=configprobes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=synthetic.clodevo.com,resources=configprobes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=synthetic.clodevo.com,resources=configprobes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigProbe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ConfigProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)
	log := ctrl.Log.WithName("controllers").WithName("ConfigProbe")

	// Fetch the ConfigProbe instance
	instance := &syntheticv1.ConfigProbe{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ConfigProbe resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ConfigProbe")
		return ctrl.Result{}, err
	}

	// Handle the configuration from instance.Spec
	// ...
	log.Info("Successfully fetched ConfigProbe", "ConfigProbe", instance)

	// Serialize the instance.Spec to a file
	fileContent, err := json.Marshal(instance.Spec)
	if err != nil {
		log.Error(err, "Failed to serialize ConfigProbe spec")
		return ctrl.Result{}, err
	}

	// Define the directory and file paths based on the emptyDir volume
	dirPath := fmt.Sprintf("/data/%s", instance.Namespace)
	filePath := fmt.Sprintf("%s/%s.json", dirPath, instance.Name)

	// Ensure the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
	}

	// Write the serialized spec to the file
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		log.Error(err, "Failed to write ConfigProbe spec to file")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigProbeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syntheticv1.ConfigProbe{}).
		Complete(r)
}
