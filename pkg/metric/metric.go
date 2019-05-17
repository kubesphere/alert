package metric

import ()

type MetricParam struct {
	RsTypeName       string              `json:"rs_type_name"`
	RsTypeParam      string              `json:"rs_type_param"`
	RsFilterName     string              `json:"rs_filter_name"`
	RsFilterParam    string              `json:"rs_filter_param"`
	ExtraQueryParams string              `json:"extra_query_params"`
	Metrics          []string            `json:"metrics"`
	MetricToRule     map[string][]string `json:"metric_to_rule"`
}

type TV struct {
	T int64  `json:"time"`
	V string `json:"value"`
}

type ResourceMetrics struct {
	RuleId         string
	MetricName     string
	ResourceMetric map[string][]TV
}
