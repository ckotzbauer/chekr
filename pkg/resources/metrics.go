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

func CPUCores(namespace, pod string) string {
	return fmt.Sprintf("{__name__=~\"kube_pod_container_resource_requests_cpu_cores|kube_pod_container_resource_limits_cpu_cores|node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate\", namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func Memory(namespace, pod string) string {
	return fmt.Sprintf("{__name__=~\"kube_pod_container_resource_requests_cpu_cores|kube_pod_container_resource_limits_cpu_cores|node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate|container_memory_working_set_bytes|kube_pod_container_resource_requests_memory_bytes|kube_pod_container_resource_limits_memory_bytes\", namespace=\"%v\", pod=\"%v\", container!=\"\", container!=\"POD\"}", namespace, pod)
}

/*func CpuCoresRequests(namespace, pod string) string {
	return fmt.Sprintf("kube_pod_container_resource_requests_cpu_cores{namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func CpuCoresLimits(namespace, pod string) string {
	return fmt.Sprintf("kube_pod_container_resource_limits_cpu_cores{namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func CpuCoresUsage(namespace, pod string) string {
	return fmt.Sprintf("node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func MemoryRequests(namespace, pod string) string {
	return fmt.Sprintf("kube_pod_container_resource_requests_memory_bytes{namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func MemoryLimits(namespace, pod string) string {
	return fmt.Sprintf("kube_pod_container_resource_limits_memory_bytes{namespace=\"%v\", pod=\"%v\"}", namespace, pod)
}

func MemoryUsage(namespace, pod string) string {
	return fmt.Sprintf("container_memory_working_set_bytes{namespace=\"%v\", pod=\"%v\", container!=\"\", container!=\"POD\", id=~\"/system.slice.*\"}", namespace, pod)
}*/
