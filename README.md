# namespacemonitor

Monitor the resources (CPU, memory) usage of Pod and containerin the namespace on Kubernetes.

## Description

This is kubebuilder operator to provide the monitor of the Pod and container CPU and memory resource images in the specific namespace on Kubernetes environment.

## Getting Started

### Prerequisites

Based on the local development environment, such a high version number is not actually required.

- go version v1.23.0+
- podman version 5.2.2+.
- kubectl version v1.27.0+.
- Access to a Kubernetes v1.27.0+ cluster.

### To Deploy Operator on the cluster

Introduce some different deployment methods

#### 1. Debug testing on the local with source code

**Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects:**

```sh
make manifests
```

**Install the CRDs into the cluster:**

```sh
make install
```

**Run a Operator Controller from host:**

```sh
make run
```

**NOTE:** Use the "Ctrl+C" to cancel it.

#### 2.Deploy with source code

**Build your image to the location specified by `IMG`:**

```sh
make docker-build IMG=<some-registry>/namespacemonitor:tag
```

**NOTE:** The container image building tool (podman or docker ) defined on the Makefile

**Push the image to the local registry:**

For example:

```sh
make docker-push IMG=<some-registry>/namespacemonitor:tag
```

**NOTE:** It's based on your local registry.

**Install the CRDs into the cluster:**

```sh
make manifests && make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/namespacemonitor:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

**UnDeploy:**

```sh
make undeploy
```

#### 3.Deploy with helm chart

**Create namespace on the Kuberents environment:**

```sh
kubectl create namespace namespacemonitor-system
```

**Download and onborading the container image to your kubernetes:**

`Image located on the image/*.tar`

**Download the helm chart and upate the values-override.yaml:**

```
Helm chart located on the helm/
```

The the image registry information to values-override.yaml

```
# Default values for ..
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

# This sets the container image more information can be found here: https://kubernetes.io/docs/concepts/containers/images/
image:
  repository: localhost/namespacemonitor-operator
  # This sets the pull policy for images.
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v0.0.5"
```

**Deploy the operator:**

```sh
helm install namespacemonitor-operator ./ --namespace namespacemonitor-system --debug
```

**UnDeploy the operator:**

```sh
helm uninstall namespacemonitor-operator --namespace namespacemonitor-system --debug
```

### **To Deploy CR**

You can apply the samples (examples) from the config/sample/monitoring_v1_namespacemonitor.yaml:

```sh
kubectl apply -f config/samples/monitoring_v1_namespacemonitor.yaml
```

> **NOTE**: Ensure that the samples has real values to test it out.


## Contributing

Welcom Issues and tickets.

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)


## License

MIT
