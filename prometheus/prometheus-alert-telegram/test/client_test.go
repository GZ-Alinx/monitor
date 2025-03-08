package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendAlert(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var alertData struct {
			Alerts []struct {
				Status      string            `json:"status"`
				Labels      map[string]string `json:"labels"`
				Annotations map[string]string `json:"annotations"`
			} `json:"alerts"`
		}

		err := json.NewDecoder(r.Body).Decode(&alertData)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(alertData.Alerts))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alert received and processed"))
	}))
	defer ts.Close()

	// 模拟告警数据
	alertData := struct {
		Alerts []struct {
			Status      string            `json:"status"`
			Labels      map[string]string `json:"labels"`
			Annotations map[string]string `json:"annotations"`
		} `json:"alerts"`
	}{
		Alerts: []struct {
			Status      string            `json:"status"`
			Labels      map[string]string `json:"labels"`
			Annotations map[string]string `json:"annotations"`
		}{
			{
				Status: "resolved",
				Labels: map[string]string{
					"alertname":   "HighCPUUsage",
					"severity":    "开发验证",
					"environment": "测试环境",
				},
				Annotations: map[string]string{
					"summary":     "CPU usage is above threshold",
					"description": "The CPU usage is above the defined threshold",
				},
			},
		},
	}

	// 序列化告警数据
	body, err := json.Marshal(alertData)
	assert.NoError(t, err)

	// 创建请求
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 验证响应
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
