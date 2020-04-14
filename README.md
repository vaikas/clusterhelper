# ClusterHelper

[![GoDoc](https://godoc.org/knative.dev/sample-controller?status.svg)](https://godoc.org/knative.dev/sample-controller)
[![Go Report Card](https://goreportcard.com/badge/knative/sample-controller)](https://goreportcard.com/report/knative/sample-controller)

## Motivation

In some clusters by default there's some manual mungering that needs to be done
for each namespace. For example, laying out a sensible RoleBinding for the
'default' service account or installing image pull secrets. This is a very
purpose built robot that handles a few things for you:

1. Creates a `RoleBinding` for the default `Service Account` for the selected
   (more about that later). 
1. Copies the specified `Secret` to the `Namespace`
1. Patches the above mentioned `Secret` to imagePullSecrets to `default`
   `Service Account` 

So, the idea is that you'd run this Controller and specify the policy (which
namespaces to inject), what to inject them with and hopefully your default
service account is now set up with some sane defaults. This work was motivated
by running on some clusters where for each namespace, I ended up doing the same
things over and over again. There are probably better ways to do this, but this
is what I ended up doing.

## Configuration

In the config/controller.yaml file there are 4 environment flags you want to
configure before firing things up (since they are env variables, you should set
those up before creating the controller).

1. CLUSTER_ROLE - Defines which cluster role to create a `Role Binding` for the
          default Service Account. For example, if there are `Pod Security
          Policies` that you'd like to have each default service account to
          conform to, you could use it here.
1. SOURCE_SECRET_NAMESPACE - Defines which namespaces holds to `Image Pull
          Secret`s that should be used. If you're using a private registry, you
          could stash those creds into a single namespace and distribute them
          from there.
1. SOURCE_SECRET_NAME - Holds the name of the actual `Image Pull Secret` that
          should be used.
1. INJECTION_DEFAULT - Defines whether the injection is on by default (all
          namespaces will get modified (unless labeled as not being injected)),
          or off by default (you have to label a namespace for it to get injected).

## Installation

Once you've configured the above settings, you might want to choose which
namespace you want the code to run in. Set the resource namespaces in the
config/* directory to your liking and then just fire it up.

```shell
ko apply -f ./config/
```

## What it does

So, overall the robot is pretty straightforward, it watches for all the
namespaces, and decides if it should inject/modify resources as described
above. If the particular namespace is to be injected, it will:

1. Create a `RoleBinding` in that namespace where the roleRef is the
`ClusterRole` you specified in the `CLUSTER_ROLE` env variable above and subject
is the `default` `Service Account` in that `Namespace`.
1. Copy the Secret from the specified
   `SOURCE_SECRET_NAMESPACE`/`SOURCE_SECRET_NAME` to the namespace.
1. Patch the `default` `Service Account` with the newly created secret as an
   `imagePullSecret`.
