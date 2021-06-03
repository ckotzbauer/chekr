package resources

import "fmt"

const (
	MemoryUsageMetric    = "container_memory_working_set_bytes"
	MemoryRequestsMetric = "kube_pod_container_resource_requests_memory_bytes"
	MemoryLimitsMetric   = "kube_pod_container_resource_limits_memory_bytes"
	CPUUsageMetric       = "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate"
	CPURequestsMetric    = "kube_pod_container_resource_requests_cpu_cores"
	CPULimitsMetric      = "kube_pod_container_resource_limits_cpu_cores"
)

func Metrics(namespace, pod string) string {
	return fmt.Sprintf("{__name__=~\"kube_pod_container_resource_requests_cpu_cores|kube_pod_container_resource_limits_cpu_cores|node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate|container_memory_working_set_bytes|kube_pod_container_resource_requests_memory_bytes|kube_pod_container_resource_limits_memory_bytes\", namespace=\"%v\", pod=\"%v\", container!=\"\", container!=\"POD\"}", namespace, pod)
}
