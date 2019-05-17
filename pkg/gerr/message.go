// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import "fmt"

type ErrorMessage struct {
	Name string
	en   string
	zhCN string
}

func (em ErrorMessage) Message(locale string, err error, a ...interface{}) string {
	format := ""
	switch locale {
	case En:
		format = em.en
	case ZhCN:
		if len(em.zhCN) > 0 {
			format = em.zhCN
		} else {
			format = em.en
		}
	}
	if err != nil {
		return fmt.Sprintf("%s: %s", fmt.Sprintf(format, a...), err.Error())
	} else {
		return fmt.Sprintf(format, a...)
	}
}

var (
	ErrorMissingParameter = ErrorMessage{
		Name: "missing_parameter",
		en:   "missing parameter [%s]",
		zhCN: "缺少参数[%s]",
	}

	ErrorStringLengthExceed = ErrorMessage{
		Name: "string_length_exceed",
		en:   "[%s] Length Exceed [%d]",
		zhCN: "[%s]字符串长度超过[%d]限制",
	}
	ErrorIllegalTimeFormat = ErrorMessage{
		Name: "illegal_time_format",
		en:   "illegal Time format [%s]",
		zhCN: "非法的时间格式[%s]",
	}
	ErrorDescribeResourcesFailed = ErrorMessage{
		Name: "describe_resources_failed",
		en:   "describe resources failed",
		zhCN: "获取资源失败",
	}
	ErrorCreateResourcesFailed = ErrorMessage{
		Name: "create_resources_failed",
		en:   "create resources failed",
		zhCN: "创建资源失败",
	}
	ErrorUpdateResourceFailed = ErrorMessage{
		Name: "update_resource_failed",
		en:   "update resource [%s] failed",
		zhCN: "更新资源[%s]失败",
	}
	ErrorUnsupportedParameterValue = ErrorMessage{
		Name: "unsupported_parameter_value",
		en:   "unsupported parameter [%s] value [%s]",
		zhCN: "参数[%s]不支持值[%s]",
	}
	ErrorInternalError = ErrorMessage{
		Name: "internal_error",
		en:   "internal error",
		zhCN: "内部错误",
	}
	ErrorValidateFailed = ErrorMessage{
		Name: "validate_failed",
		en:   "validate failed",
		zhCN: "校验失败",
	}
	ErrorDeleteResourceFailed = ErrorMessage{
		Name: "delete_resource_failed",
		en:   "delete resource [%s] failed",
		zhCN: "删除资源[%s]失败",
	}
)
