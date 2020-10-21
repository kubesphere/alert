update metric set metric_name="workload_cpu_usage" where metric_name="workload_pod_cpu_usage";
update metric set metric_name="workload_memory_usage" where metric_name="workload_pod_memory_usage";
update metric set metric_name="workload_memory_usage_wo_cache" where metric_name="workload_pod_memory_usage_wo_cache";
update metric set metric_name="workload_net_bytes_transmitted" where metric_name="workload_pod_net_bytes_transmitted";
update metric set metric_name="workload_net_bytes_received" where metric_name="workload_pod_net_bytes_received";
