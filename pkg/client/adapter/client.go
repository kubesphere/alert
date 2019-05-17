package adapter

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/logger"
)

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   10000,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

const (
	DefaultScheme = "http"
)

func SendMetricRequest(metricParam string) string {
	cfg := config.GetInstance()
	url := fmt.Sprintf("http://127.0.0.1:%s/api/v1/metric", cfg.App.AdapterPort)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(nil, "SendMetricRequest NewRequest error:", err)
		return ""
	}

	params := request.URL.Query()
	params.Add("metric_param", metricParam)
	request.URL.RawQuery = params.Encode()

	logger.Debug(nil, "SendMetricRequest %s", request.URL.RawQuery)

	response, err := client.Do(request)
	if err != nil {
		logger.Error(nil, "SendMetricRequest get error:", err)
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			logger.Error(nil, "SendMetricRequest read error:", err.Error())
		}

		return string(contents)
	}
	return ""
}

func SendEmailRequest(notificationParam string) string {
	cfg := config.GetInstance()
	url := fmt.Sprintf("http://127.0.0.1:%s/api/v1/email", cfg.App.AdapterPort)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(nil, "SendEmailRequest NewRequest error:", err)
		return ""
	}

	params := request.URL.Query()
	params.Add("notification_param", notificationParam)
	request.URL.RawQuery = params.Encode()

	logger.Debug(nil, "SendEmailRequest %s", request.URL.RawQuery)

	response, err := client.Do(request)
	if err != nil {
		logger.Error(nil, "SendEmailRequest get error:", err)
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			logger.Error(nil, "SendEmailRequest read error:", err.Error())
		}

		return string(contents)
	}
	return ""
}
