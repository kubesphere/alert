package notification

import ()

type NotificationParam struct {
	ResourceName   string `json:"resource_name"`
	RuleName       string `json:"rule_name"`
	CumulatedCount uint32 `json:"cumulated_count"`
	FirstTime      string `json:"first_time"`
	LastTime       string `json:"last_time"`
	LastValue      string `json:"last_value"`
}

type Email struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
