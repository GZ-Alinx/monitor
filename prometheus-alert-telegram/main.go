package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// TelegramMessage 用于发送到 Telegram 的消息结构
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// Config 包含应用程序配置
type Config struct {
	BotToken string
	ChatID   string
	Port     string
}

// Alert 结构体用于解析 Prometheus 告警的 JSON 数据
type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// logger 全局日志对象
var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("无法初始化日志: %v", err)
	}
}

// loadConfig 从环境变量加载配置
func loadConfig() (*Config, error) {
	config := &Config{
		BotToken: "7996277665:AAFWXkzS8iyUNJyiTFxoKTSVzOtgP6EGLQE",
		ChatID:   "-4508961268",
		Port:     ":9080",
	}

	return config, nil
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// webhookHandler 处理 Prometheus 发送的告警
func webhookHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var alerts struct {
			Alerts []Alert `json:"alerts"`
		}

		if err := json.NewDecoder(r.Body).Decode(&alerts); err != nil {
			logger.Error("Failed to decode JSON", zap.Error(err))
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		for _, alert := range alerts.Alerts {
			message := formatAlertMessage(alert)
			if err := sendTelegramMessage(config.BotToken, config.ChatID, message); err != nil {
				logger.Error("Failed to send alert to Telegram",
					zap.Error(err),
					zap.String("chat_id", config.ChatID))
				http.Error(w, "Failed to send alert", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alert received and processed"))
	}
}

// formatAlertMessage 格式化告警消息为字符串
func formatAlertMessage(alert Alert) string {
	status := "🚨触发"
	if alert.Status == "resolved" {
		status = "✅恢复"
	}

	environment := alert.Labels["environment"]
	if environment == "" {
		environment = "未知环境"
	}

	return fmt.Sprintf(
		"*%s告警 - %s*\n*告警名称*: %s\n*严重级别*: %s\n*摘要*: %s\n*描述*: %s\n",
		status,
		environment,
		alert.Labels["alertname"],
		alert.Labels["severity"],
		alert.Annotations["summary"],
		alert.Annotations["description"],
	)
}

// sendTelegramMessage 发送消息到 Telegram
func sendTelegramMessage(token, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?parse_mode=Markdown", token)

	telegramMessage := TelegramMessage{
		ChatID: chatID,
		Text:   message,
	}

	body, err := json.Marshal(telegramMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal Telegram message: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request to Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API returned status code %d: %s", resp.StatusCode, string(respBody))
	}

	logger.Info("Alert sent to Telegram", zap.String("chat_id", chatID))
	return nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	http.HandleFunc("/alert", webhookHandler(config))

	logger.Info("Starting server",
		zap.String("port", config.Port),
		zap.String("chat_id", config.ChatID))

	server := &http.Server{
		Addr:         config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
