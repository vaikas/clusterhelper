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

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"knative.dev/pkg/kmeta"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	"github.com/vaikas/clusterhelper/pkg/reconciler/clusterhelper/resources"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/namespace"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/secret"
	sainformer "knative.dev/pkg/client/injection/kube/informers/core/v1/serviceaccount"
	rbacinformer "knative.dev/pkg/client/injection/kube/informers/rbac/v1/rolebinding"
	namespacereconciler "knative.dev/pkg/client/injection/kube/reconciler/core/v1/namespace"
)

const (
	// GroupName is the our group name.
	GroupName = "clusterhelper.internal.vaikas.dev"

	// This is what we label our resources with for easy filtering.
	ResourceLabelKey = GroupName + "/clusterhelperkey"

	// This is what we set the value to for easy filtering.
	ResourceLabelValue = "true"
)

type envConfig struct {
	InjectionDefault      bool   `envconfig:"CLUSTER_HELPER_INJECTION_DEFAULT" default:"true"`
	ClusterRole           string `envconfig:"CLUSTER_ROLE" required:"true"`
	SourceSecretNamespace string `envconfig:"SOURCE_SECRET_NAMESPACE" required:"true"`
	SourceSecretName      string `envconfig:"SOURCE_SECRET_NAME" required:"true"`
}

func onByDefault(labels map[string]string) bool {
	return labels[resources.InjectionLabelKey] == resources.InjectionDisabledLabelValue
}

func offByDefault(labels map[string]string) bool {
	return labels[resources.InjectionLabelKey] != resources.InjectionEnabledLabelValue
}

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)

	namespaceInformer := namespace.Get(ctx)
	secretInformer := secret.Get(ctx)
	serviceaccountInformer := sainformer.Get(ctx)
	rolebindingInformer := rbacinformer.Get(ctx)

	var filter labelFilter

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logging.FromContext(ctx).Fatalf("clusterhelper was unable to process environment: %v", err)
	} else if env.InjectionDefault {
		filter = onByDefault
	} else {
		filter = offByDefault
	}

	r := &Reconciler{
		clusterRole:           env.ClusterRole,
		sourceSecretNamespace: env.SourceSecretNamespace,
		sourceSecretName:      env.SourceSecretName,
		filter:                filter,
		namespaceLister:       namespaceInformer.Lister(),
		kubeclient:            kubeclient.Get(ctx),
		rbacLister:            rolebindingInformer.Lister(),
		saLister:              serviceaccountInformer.Lister(),
		secretLister:          secretInformer.Lister(),
	}
	impl := namespacereconciler.NewImpl(ctx, r)

	logger.Info("Setting up event handlers.")

	namespaceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Note that we set up the informers to trigger a namespace reconcile when any
	// of our resources change. We do this by using a label so that we can easily
	// filter on them.
	secretInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: pkgreconciler.LabelFilterFunc(ResourceLabelKey, ResourceLabelValue, false),
		Handler:    controller.HandleAll(EnqueueNamespaceOf(ctx, impl)),
	})
	// Only care about the default service account.
	serviceaccountInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: pkgreconciler.NameFilterFunc("default"),
		Handler:    controller.HandleAll(EnqueueNamespaceOf(ctx, impl)),
	})
	rolebindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: pkgreconciler.LabelFilterFunc(ResourceLabelKey, ResourceLabelValue, false),
		Handler:    controller.HandleAll(EnqueueNamespaceOf(ctx, impl)),
	})

	return impl
}

func EnqueueNamespaceOf(ctx context.Context, impl *controller.Impl) func(obj interface{}) {
	return func(obj interface{}) {
		object, err := kmeta.DeletionHandlingAccessor(obj)
		if err != nil {
			logging.FromContext(ctx).Errorw("Enqueue", zap.Error(err))
			return
		}
		impl.EnqueueKey(types.NamespacedName{Name: object.GetNamespace()})
	}
}
