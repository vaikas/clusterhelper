/*
Copyright 2019 The Knative Authors

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

package clusterhelper

import (
	"context"
	"fmt"

	"github.com/vaikas/clusterhelper/pkg/reconciler/clusterhelper/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	rbacv1listers "k8s.io/client-go/listers/rbac/v1"
	"knative.dev/pkg/apis/duck"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	namespacereconciler "knative.dev/pkg/client/injection/kube/reconciler/core/v1/namespace"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	pkgreconciler "knative.dev/pkg/reconciler"
)

const (
	resourceName = "cluster-helper"
)

// newReconciledNormal makes a new reconciler event with event type Normal, and
// reason ClusterHelperNamespaceReconciled.
func newReconciledNormal(namespace, name string) pkgreconciler.Event {
	return reconciler.NewEvent(corev1.EventTypeNormal, "ClusterHelperNamespaceReconciled", "Namespace reconciled by ClusterHelper: \"%s/%s\"", namespace, name)
}

type labelFilter func(labels map[string]string) bool

// Reconciler implements ReconcileKind for Namespaces.
type Reconciler struct {
	kubeclient kubernetes.Interface

	sourceSecretNamespace string
	sourceSecretName      string
	clusterRole           string

	namespaceLister corev1listers.NamespaceLister
	rbacLister      rbacv1listers.RoleBindingLister
	secretLister    corev1listers.SecretLister
	saLister        corev1listers.ServiceAccountLister

	filter labelFilter
}

// Check that our Reconciler implements Interface
var _ namespacereconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, ns *corev1.Namespace) pkgreconciler.Event {
	if r.filter(ns.Labels) {
		logging.FromContext(ctx).Debug("Not reconciling Namespace: %q", ns.Name)
		return nil
	}

	if ns.GetDeletionTimestamp() != nil {
		// Check for a DeletionTimestamp.  If present, elide the normal reconcile logic.
		// When a controller needs finalizer handling, it would go here.
		return nil
	}

	// Only create this if clusterRole has been specified
	if r.clusterRole != "" {
		if err := r.reconcileRoleBinding(ctx, ns.Name); err != nil {
			return err
		}
	}

	if r.sourceSecretNamespace != "" && r.sourceSecretName != "" {
		if err := r.reconcileSecret(ctx, ns.Name); err != nil {
			return err
		}

		if err := r.reconcileServiceAccount(ctx, ns.Name); err != nil {
			return err
		}
	}

	return newReconciledNormal(ns.Namespace, ns.Name)
}

func (r *Reconciler) reconcileRoleBinding(ctx context.Context, namespace string) error {
	name := resourceName

	roleBinding, err := r.rbacLister.RoleBindings(namespace).Get(resourceName)
	if apierrs.IsNotFound(err) {
		roleBinding = resources.MakeRoleBinding(ctx, name, namespace, r.clusterRole)
		roleBinding, err = r.kubeclient.RbacV1().RoleBindings(namespace).Create(roleBinding)
		if err != nil {
			return fmt.Errorf("failed to create rolebinding %q: %w", name, err)
		}
		logging.FromContext(ctx).Infof("Created rolebinding %q", name)
	} else if err != nil {
		return fmt.Errorf("failed to get rolebinding %q: %w", name, err)
	}
	return nil
}

func (r *Reconciler) reconcileSecret(ctx context.Context, namespace string) error {
	name := resourceName
	_, err := r.secretLister.Secrets(namespace).Get(resourceName)
	if apierrs.IsNotFound(err) {
		secret, err := r.secretLister.Secrets(r.sourceSecretNamespace).Get(r.sourceSecretName)
		newSecret := secret.DeepCopy()
		newSecret.ObjectMeta = metav1.ObjectMeta{Namespace: namespace, Name: name}
		secret, err = r.kubeclient.CoreV1().Secrets(namespace).Create(newSecret)
		if err != nil {
			return fmt.Errorf("failed to create secret %q: %w", name, err)
		}
		logging.FromContext(ctx).Infof("Created secret %q", name)
	} else if err != nil {
		return fmt.Errorf("failed to get secret %q: %w", name, err)
	}

	return nil
}

func (r *Reconciler) reconcileServiceAccount(ctx context.Context, namespace string) error {
	logger := logging.FromContext(ctx)

	sa, err := r.saLister.ServiceAccounts(namespace).Get("default")
	newSa := sa.DeepCopy()
	if err != nil {
		return fmt.Errorf("failed to get default service account: %w", err)
	}

	for _, is := range newSa.ImagePullSecrets {
		if is.Name == resourceName {
			return nil
		}
	}

	newSa.ImagePullSecrets = append(newSa.ImagePullSecrets, corev1.LocalObjectReference{Name: resourceName})

	patch, err := duck.CreateMergePatch(sa, newSa)
	if err != nil {
		logger.Warnf("Failed to create patch for \"%s/%s\" : %v", namespace, "default", err)
		return err
	}
	logger.Infof("Patched \"%s/%s\": %q", namespace, "default", string(patch))
	// If there is nothing to patch, we are good, just return.
	// Empty patch is {}, hence we check for that.
	if len(patch) <= 2 {
		return nil
	}

	patched, err := r.kubeclient.CoreV1().ServiceAccounts(namespace).Patch("default", types.MergePatchType, patch)
	if err != nil {
		logger.Warnf("Failed to patch \"%s/%s\" : %v", namespace, "default", err)
		return err
	}
	logger.Infof("new Service Account : %+v", patched)

	return nil
}
