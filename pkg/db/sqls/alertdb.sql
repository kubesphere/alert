SET SESSION FOREIGN_KEY_CHECKS=0;

/* Drop Tables */

DROP TABLE IF EXISTS action;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS alert;
DROP TABLE IF EXISTS rule;
DROP TABLE IF EXISTS metric;
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
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
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
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	policy_id varchar(50) NOT NULL,
	rs_filter_id varchar(50) NOT NULL,
	executor_id varchar(50),
	PRIMARY KEY (alert_id)
);


CREATE TABLE comment
(
	comment_id varchar(50) NOT NULL,
	addresser varchar(50) NOT NULL,
	content varchar(255) NOT NULL,
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	history_id varchar(50) NOT NULL,
	PRIMARY KEY (comment_id)
);


CREATE TABLE history
(
	history_id varchar(50) NOT NULL,
	history_name varchar(50) NOT NULL,
	-- triggered resumed sent_success sent_failed commentd
	event varchar(50) NOT NULL COMMENT 'triggered resumed sent_success sent_failed commentd',
	content text,
	notification_id varchar(50),
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	alert_id varchar(50) NOT NULL,
	rule_id varchar(50) NOT NULL,
	resource_name varchar(300),
	PRIMARY KEY (history_id)
);


CREATE TABLE metric
(
	metric_id varchar(50) NOT NULL,
	metric_name varchar(100) DEFAULT '' NOT NULL,
	metric_param text NOT NULL,
	-- active inactive
	status varchar(20) DEFAULT 'active' NOT NULL COMMENT 'active inactive',
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (metric_id)
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
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (policy_id)
);


CREATE TABLE resource_filter
(
	rs_filter_id varchar(50) NOT NULL,
	rs_filter_name varchar(50) NOT NULL,
	rs_filter_param text,
	-- active,disabled,deleted
	status varchar(50) DEFAULT 'active' NOT NULL COMMENT 'active,disabled,deleted',
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	rs_type_id varchar(50) NOT NULL,
	PRIMARY KEY (rs_filter_id)
);


CREATE TABLE resource_type
(
	rs_type_id varchar(50) NOT NULL,
	rs_type_name varchar(50) DEFAULT '' NOT NULL,
	rs_type_param text,
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
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
	-- datetime(3)
	create_time datetime(3) COMMENT 'datetime(3)',
	-- datetime(3)
	update_time datetime(3) COMMENT 'datetime(3)',
	policy_id varchar(50) NOT NULL,
	metric_id varchar(50) NOT NULL,
	PRIMARY KEY (rule_id)
);



/* Create Foreign Keys */

ALTER TABLE history
	ADD FOREIGN KEY (alert_id)
	REFERENCES alert (alert_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE comment
	ADD FOREIGN KEY (history_id)
	REFERENCES history (history_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE rule
	ADD FOREIGN KEY (metric_id)
	REFERENCES metric (metric_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE action
	ADD FOREIGN KEY (policy_id)
	REFERENCES policy (policy_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE alert
	ADD FOREIGN KEY (policy_id)
	REFERENCES policy (policy_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE rule
	ADD FOREIGN KEY (policy_id)
	REFERENCES policy (policy_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE alert
	ADD FOREIGN KEY (rs_filter_id)
	REFERENCES resource_filter (rs_filter_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE metric
	ADD FOREIGN KEY (rs_type_id)
	REFERENCES resource_type (rs_type_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE policy
	ADD FOREIGN KEY (rs_type_id)
	REFERENCES resource_type (rs_type_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE resource_filter
	ADD FOREIGN KEY (rs_type_id)
	REFERENCES resource_type (rs_type_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE history
	ADD FOREIGN KEY (rule_id)
	REFERENCES rule (rule_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;



