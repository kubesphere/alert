/* Drop Tables */

DROP TABLE IF EXISTS action;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS alert;
DROP TABLE IF EXISTS executor;
DROP TABLE IF EXISTS rule;
DROP TABLE IF EXISTS metric;
DROP TABLE IF EXISTS nf_address_list;
DROP TABLE IF EXISTS policy;
DROP TABLE IF EXISTS resource_filter;
DROP TABLE IF EXISTS resource_type;




/* Create Tables */

CREATE TABLE action
(
	action_id varchar(50) NOT NULL,
	action_name varchar(50) NOT NULL,
	-- if alert change to this status, trigger action，监控状态变为正常，告警，监控数据不足时触发行为。
	-- Turn2Normal
	-- Turn2Alarm
	-- Both
	trigger_status varchar(50) DEFAULT 'Turn2Alarm' NOT NULL COMMENT 'if alert change to this status, trigger action，监控状态变为正常，告警，监控数据不足时触发行为。
Turn2Normal
Turn2Alarm
Both',
	-- notification_list_id
	trigger_action varchar(255) COMMENT 'notification_list_id',
	created_at timestamp,
	updated_at timestamp,
	policy_id varchar(50) NOT NULL,
	nf_address_list_id varchar(50) NOT NULL,
	PRIMARY KEY (action_id)
);


CREATE TABLE alert
(
	alert_id varchar(50) NOT NULL,
	alert_name varchar(100) NOT NULL,
	disabled boolean DEFAULT false NOT NULL,
	-- running,adding,deleting,updating,migrating
	running_status varchar(50) NOT NULL COMMENT 'running,adding,deleting,updating,migrating',
	-- map[rule_id] {severity cumulated_send_count next_resend_interval next_sendable_time aggregated_alert}
	alert_status text NOT NULL COMMENT 'map[rule_id] {severity cumulated_send_count next_resend_interval next_sendable_time aggregated_alert}',
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	policy_id varchar(50) NOT NULL,
	rs_filter_id varchar(50) NOT NULL,
	executor_id varchar(50) NOT NULL,
	PRIMARY KEY (alert_id)
);


CREATE TABLE comment
(
	comment_id varchar(50) NOT NULL,
	addresser varchar(50) NOT NULL,
	content varchar(100) NOT NULL,
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	history_id varchar(50) NOT NULL,
	PRIMARY KEY (comment_id)
);


CREATE TABLE executor
(
	executor_id varchar(50) NOT NULL,
	executor_name varchar(50),
	-- running,down,deleted
	--  
	status varchar(20) DEFAULT 'running' NOT NULL COMMENT 'running,down,deleted
 ',
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	PRIMARY KEY (executor_id)
);


CREATE TABLE history
(
	history_id varchar(50) NOT NULL,
	history_name varchar(50) NOT NULL,
	-- triggered resumed sent_success sent_failed commentd
	event varchar(50) NOT NULL COMMENT 'triggered resumed sent_success sent_failed commentd',
	content text,
	notification_id varchar(50),
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	alert_id varchar(50) NOT NULL,
	rule_id varchar(50) NOT NULL,
	resource_name varchar(50),
	PRIMARY KEY (history_id)
);


CREATE TABLE metric
(
	metric_id varchar(50) NOT NULL,
	metric_name varchar(50) DEFAULT '' NOT NULL,
	metric_param text NOT NULL,
	-- active inactive
	status varchar(20) DEFAULT 'active' NOT NULL COMMENT 'active inactive',
	created_at timestamp DEFAULT NOW(),
	updated_at timestamp,
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (metric_id)
);


CREATE TABLE nf_address_list
(
	nf_address_list_id varchar(50) NOT NULL,
	PRIMARY KEY (nf_address_list_id)
);


CREATE TABLE policy
(
	policy_id varchar(50) NOT NULL,
	policy_name varchar(50) NOT NULL,
	policy_description varchar(255),
	-- map[severity] {repeat_type repeat_interval_initvalue max_send_count}
	policy_config text COMMENT 'map[severity] {repeat_type repeat_interval_initvalue max_send_count}',
	creator varchar(50),
	available_start_time varchar(8),
	available_end_time varchar(8),
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (policy_id)
);


CREATE TABLE resource_filter
(
	rs_filter_id varchar(50) NOT NULL,
	rs_filter_name varchar(50) NOT NULL,
	rs_filter_uri text,
	-- active,disabled,deleted
	status varchar(50) DEFAULT 'active' NOT NULL COMMENT 'active,disabled,deleted',
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (rs_filter_id)
);


CREATE TABLE resource_type
(
	rs_type_id varchar(50) NOT NULL,
	rs_type_name varchar(50) DEFAULT '' NOT NULL,
	rs_uri_tmpl text,
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	PRIMARY KEY (rs_type_id)
);


CREATE TABLE rule
(
	rule_id varchar(50) NOT NULL,
	rule_name varchar(50) NOT NULL,
	disabled boolean DEFAULT false NOT NULL,
	-- unit：minute
	monitor_periods int DEFAULT 5 NOT NULL COMMENT 'unit：minute',
	-- critical,major,minor
	severity varchar(20) DEFAULT 'minor' NOT NULL COMMENT 'critical,major,minor',
	-- 实时/均值/最小值/最大值/计数/求和
	-- instant/avg/min/max/count/sum
	metrics_type varchar(10) COMMENT '实时/均值/最小值/最大值/计数/求和
instant/avg/min/max/count/sum',
	-- = ;     != ;    >;    <;    >=;    <=; 
	condition_type varchar(10) NOT NULL COMMENT '= ;     != ;    >;    <;    >=;    <=; ',
	-- thresholds: v1|v2...
	thresholds varchar(50) NOT NULL COMMENT 'thresholds: v1|v2...',
	unit varchar(50),
	consecutive_count int DEFAULT 1 NOT NULL,
	inhibit boolean DEFAULT false NOT NULL,
	created_at timestamp,
	updated_at timestamp DEFAULT NOW(),
	policy_id varchar(50) NOT NULL,
	metric_id varchar(50) NOT NULL,
	PRIMARY KEY (rule_id)
);

INSERT INTO `metric` VALUES ('mt-021lmDGygjMY','cluster_disk_size_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-0WMZpN7MgjMY','namespace_memory_usage_wo_cache','0.000000000931322574615478515625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-1E7Kv3LA4oA7','namespace_cpu_usage','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-29oZ0oV14DnB','container_cpu_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-2M70Qwkn4DnB','pod_memory_usage','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-2NV3LVWQ5OPJ','node_disk_size_available','0.000000001','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-3pk7vRx2zAkg','node_load5','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-3Q8MMLM2XYlx','cluster_cpu_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-3WylnMWkzAkg','workload_statefulset_unavailable_replicas_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-61MWXJX04V8N','node_memory_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-69yNOQgY4V8N','cluster_pod_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-7AGJY7X2NO2q','workload_pod_memory_usage_wo_cache','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-7g83mByZNO2q','workspace_memory_usage_wo_cache','0.000000000931322574615478515625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-954yEkAyJKEm','workspace_pod_abnormal_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-9l3YpogEJKEm','workload_deployment_unavailable_replicas_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-9r8wYwOoXYlx','container_memory_utilisation_wo_cache','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-9wrDqPOwWk93','node_memory_available','0.000000000931322574615478515625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-AjJWxqBEJKEm','pod_memory_usage_wo_cache_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-AlV2RmR2JKEm','pod_net_bytes_received','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-AQ2jE8nEJKEm','workload_pod_net_bytes_transmitted','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-DxqyV8A7zAkg','container_memory_usage','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-DyGW4qPkzAkg','workload_pod_net_bytes_received','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-G9jvLRjD4oA7','container_memory_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-GEqBr1qyJKEm','node_pod_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-GNq5D4l7zAkg','container_cpu_usage','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-Gqg9EKqX4oA7','node_net_bytes_transmitted','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-j37xz5xkzAkg','namespace_pod_abnormal_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-jjWGPMlnzAkg','cluster_node_offline_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-JQGwLDvn4DnB','workload_pod_memory_usage','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-k0BrxOP2zAkg','node_disk_write_iops','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-k14oD7364DnB','workspace_memory_usage','0.000000000931322574615478515625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-KPxD6YWMgjMY','namespace_memory_usage','0.000000000931322574615478515625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-kZAM89664DnB','node_cpu_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-lB7l3Z9kzAkg','workload_daemonset_unavailable_replicas_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-lgGwk9n1XYlx','node_pod_abnormal_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-lM1LD6zkzAkg','pod_cpu_usage','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-nzLg8vkVgjMY','container_memory_usage_wo_cache','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-o6Q0mv5N5OPJ','cluster_disk_inode_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-OrgRJ9WwWk93','workspace_cpu_usage','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-OyR9nE8W4oA7','container_net_bytes_transmitted','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-oz3Lx2gR4DnB','container_net_bytes_received','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-2loEnEY6Oyzp'),('mt-OzOQ477EJKEm','namespace_resourcequota_used_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-p8Elp7K04V8N','node_disk_size_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-PKmBorl1XYlx','node_load1','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-PMLGQr3P5OPJ','pod_cpu_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-pQ6xP1mZNO2q','node_disk_read_throughput','0.000001','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-pv9X6BrY4V8N','cluster_pod_abnormal_ratio','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-QNWALxDn4DnB','pod_memory_usage_wo_cache','0.00000095367431640625','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-QvZ9wwBX4oA7','node_disk_write_throughput','0.000001','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-rE8yox22NO2q','namespace_net_bytes_received','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-rg8x8vplXYlx','pod_memory_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-rl5vGLMZNO2q','node_load15','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-rQ7Q3KvA4oA7','pod_net_bytes_transmitted','0.000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-7ENxQzjLKrpl'),('mt-RWXXoJkyJKEm','node_disk_inode_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-vnAjqwNP5OPJ','namespace_net_bytes_transmitted','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-pX1mLzBoJ3mA'),('mt-wEV54WQ19ZLR','workspace_net_bytes_transmitted','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-WJYqgPYMgjMY','workload_pod_cpu_usage','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-rBYPE5KLKrpl'),('mt-X818gr764DnB','node_disk_read_iops','1.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90'),('mt-xRA8G6D64DnB','workspace_net_bytes_received','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-yJO8173xP39G'),('mt-YPBylZ3nzAkg','cluster_memory_utilisation','100.0','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-GjnE66xplG90'),('mt-Z5zNrvq2zAkg','node_net_bytes_received','0.000000125','active','2019-03-14 00:00:00','2019-03-14 00:00:00','rst-3m8ZmxVylG90');

INSERT INTO `resource_type` VALUES ('rst-2loEnEY6Oyzp','container','/namespaces/{ns_name}/pods/{pod_name}/containers?resources_filter={container_name},/nodes/{node_id}/pods/{pod_name}/containers?resources_filter={container_name}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-3m8ZmxVylG90','node','/nodes?resources_filter={node_id}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-7ENxQzjLKrpl','pod','/namespaces/{ns_name}/pods?resources_filter={pod_name},/nodes/{node_id}/pods?resources_filter={pod_name}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-GjnE66xplG90','cluster','/clusters','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-pX1mLzBoJ3mA','namespace','/namespaces?resources_filter={ns_name}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-rBYPE5KLKrpl','workload','/namespaces/{ns_name}/workloads/{workload_kind}?resources_filter={workload_name}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-wgmgmR8o4jPL','component','/components?resources_filter={component_name}','2019-03-14 00:00:00','2019-03-14 00:00:00'),('rst-yJO8173xP39G','workspace','/workspaces?resources_filter={ws_name}','2019-03-14 00:00:00','2019-03-14 00:00:00');