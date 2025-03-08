package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Alert 结构与服务端保持一致
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// AlertsPayload 表示要发送到服务器的告警负载
type AlertsPayload struct {
	Alerts []Alert `json:"alerts"`
}

// createSampleAlert 创建一个样本告警
func createSampleAlert() Alert {
	return Alert{
		Status: "firing", // 或者 "resolved"
		Labels: map[string]string{
			"alertname":   "CPUHighUsage",
			"host":        "server-1",
			"severity":    "critical",
			"environment": "production",
		},
		Annotations: map[string]string{
			"summary":     "JeetCPU 使用率过高",
			"description": "CPU usage exceeded 95% for more than 5 minutes.",
		},
		StartsAt:     time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
		EndsAt:       time.Now().Format(time.RFC3339),
		GeneratorURL: "http://monitoring-system/alerts",
		Fingerprint:  "abcd1234efgh5678ijkl",
	}
}

// sendAlert 发送告警数据到服务器
func sendAlert(url string, payload AlertsPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	log.Printf("Alert successfully sent to %s", url)
	return nil
}

func main() {
	alert := createSampleAlert()
	payload := AlertsPayload{
		Alerts: []Alert{alert},
	}

	// 设置服务器的告警处理 URL
	serverURL := "http://localhost:9910/alert"

	// 发送告警
	err := sendAlert(serverURL, payload)
	if err != nil {
		log.Fatalf("Failed to send alert: %v", err)
	}
}
