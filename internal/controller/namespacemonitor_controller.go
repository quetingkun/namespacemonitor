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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/client-go/kubernetes"
	//"k8s.io/apimachinery/pkg/api/resource"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	monitoringv1 "resource.com/NamespaceMonitor/api/v1"
)

// NamespaceMonitorReconciler reconciles a NamespaceMonitor object
type NamespaceMonitorReconciler struct {
	client.Client
	KubeClient    kubernetes.Interface
	MetricsClient *metricsv.Clientset
	Scheme        *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.resource.com,resources=namespacemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.resource.com,resources=namespacemonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.resource.com,resources=namespacemonitors/finalizers,verbs=update
// +kubebuilder:rbac:groups=metrics.k8s.io,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NamespaceMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx).WithValues("namespaceMonitor", req.NamespacedName)
	// 创建一个新的日志记录器实例，不带上下文信息
	simpleLogger := log.Log.WithName("namespaceMonitor")

	// Fetch NamespaceMonitor instance
	namespaceMonitor := &monitoringv1.NamespaceMonitor{}
	if err := r.Get(ctx, req.NamespacedName, namespaceMonitor); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("NamespaceMonitor resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get NamespaceMonitor")
		return ctrl.Result{}, err
	}

	simpleLogger.Info("------------------------------------")
	// 打印 NamespaceMonitorSpec 中的字段
	simpleLogger.Info("NamespaceMonitor Spec",
		"Namespace", namespaceMonitor.Spec.Namespace,
		"UpdateInterval", namespaceMonitor.Spec.UpdateInterval,
	)

	// Check if the MetricsClient is initialized
	if r.MetricsClient == nil {
		logger.Error(nil, "MetricsClient is not initialized")
		return ctrl.Result{}, fmt.Errorf("MetricsClient is nil")
	}

	// Fetch metrics from the metrics server via the API
	podMetricsList, err := r.MetricsClient.MetricsV1beta1().PodMetricses(namespaceMonitor.Spec.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "Failed to fetch pod metrics")
		return ctrl.Result{}, err
	}

	// Update the NamespaceMonitor status
	namespaceMonitor.Status = monitoringv1.NamespaceMonitorStatus{
		LastUpdated: metav1.Now(),
		PodMetrics:  []monitoringv1.PodMetrics{},
	}

	// Output the metrics of Pod's CPU and memory usage
	for _, pod := range podMetricsList.Items {
		simpleLogger.Info("Pod metrics", "Pod", pod.Name)
		podMetrics := monitoringv1.PodMetrics{
			PodName:          pod.Name,
			ContainerMetrics: []monitoringv1.ContainerMetrics{},
		}

		for _, container := range pod.Containers {
			simpleLogger.Info("Container metrics", "Pod", pod.Name, "Container", container.Name, "CPU", container.Usage["cpu"], "Memory", container.Usage["memory"])

			//[debug] print container.Usage["cpu"] and container.Usage["memory"]
			//fmt.Printf("Type: %T, Value: %+v\n", container.Usage["cpu"], container.Usage["cpu"])
			//fmt.Printf("Type: %T, Value: %+v\n", container.Usage["memory"], container.Usage["memory"])

			cpuQuantity := container.Usage["cpu"]
			memoryQuantity := container.Usage["memory"]

			cpuUsage := cpuQuantity.String()
			memoryUsage := memoryQuantity.String()

			//fmt.Printf("Container: %s, CPU Usage: %s, Memory Usage: %s\n", container.Name, cpuUsage, memoryUsage)
			containerMetrics := monitoringv1.ContainerMetrics{
				ContainerName: container.Name,
				CPUUsage:      cpuUsage,
				MemoryUsage:   memoryUsage,
			}
			podMetrics.ContainerMetrics = append(podMetrics.ContainerMetrics, containerMetrics)
		}
		namespaceMonitor.Status.PodMetrics = append(namespaceMonitor.Status.PodMetrics, podMetrics)
	}

	if err := r.Status().Update(ctx, namespaceMonitor); err != nil {
		logger.Error(err, "Failed to update NamespaceMonitor status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.NamespaceMonitor{}).
		Complete(r)
}
