package controllers

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testing"
	"github.com/stretchr/testify/assert"
	syntheticv1 "github.com/lakhlaifi/configprobe-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconcile(t *testing.T) {
	// Setup the scheme
	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	syntheticv1.AddToScheme(scheme)

	// Create a mock client
	cl := fake.NewClientBuilder().WithScheme(scheme).Build()

	// Create an instance of your reconciler
	r := &ConfigProbeReconciler{
		Client: cl,
		Log:    testing.NullLogger{},
		Scheme: scheme,
	}

	// Create a ConfigProbe resource
	probe := &syntheticv1.ConfigProbe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-probe",
			Namespace: "default",
		},
	}

	// Use the mock client to create the resource
	cl.Create(context.TODO(), probe)

	// Call the Reconcile method
	_, err := r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: metav1.NamespacedName{
			Name:      "test-probe",
			Namespace: "default",
		},
	})
	assert.NoError(t, err)

	// Check if a ConfigMap was created
	cm := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), metav1.NamespacedName{
		Name:      "test-probe-configmap",
		Namespace: "default",
	}, cm)
	assert.NoError(t, err)
}
