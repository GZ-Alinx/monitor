package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 定义常量用于存储所有的 BotToken 和 ChatID
const (
	// JeetBotToken = "7996277665:AAFWXkzS8iyUNJyiTFxoKTSVzOtgP6EGLQE"
	// JeetChatID   = "-4508961268"

	JeetBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	JeetChatID   = "-1002428509416"

	Yk777BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	Yk777ChatID   = "-1002428509416"

	U77BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	U77ChatID   = "-1002428509416"

	AsgameBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	AsgameChatID   = "-1002428509416"

	Ez777BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	Ez777ChatID   = "-1002428509416"

	SlotinrBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	SlotinrChatID   = "-1002428509416"

	RR9BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	RR9ChatID   = "-1002428509416"

	A777winBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	A777winChatID   = "-1002428509416"

	LuckbetBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	LuckbetChatID   = "-1002428509416"

	Yy6BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	Yy6ChatID   = "-1002428509416"

	A01GameBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	A01GameChatID   = "-1002428509416"

	P77BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	P77ChatID   = "-1002428509416"

	WininrBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	WininrChatID   = "-1002428509416"

	Ww5BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	Ww5ChatID   = "-1002428509416"

	A77betBotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	A77betChatID   = "-1002428509416"

	Dd1BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	Dd1ChatID   = "-1002428509416"

	A5222BotToken = "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"
	A5222ChatID   = "-1002428509416"
)

// 初始化所有的 ChatConfig 配置
var (
	Jeet    = ChatConfig{BotToken: JeetBotToken, ChatID: JeetChatID}
	Yk777   = ChatConfig{BotToken: Yk777BotToken, ChatID: Yk777ChatID}
	U77     = ChatConfig{BotToken: U77BotToken, ChatID: U77ChatID}
	Asgame  = ChatConfig{BotToken: AsgameBotToken, ChatID: AsgameChatID}
	Ez777   = ChatConfig{BotToken: Ez777BotToken, ChatID: Ez777ChatID}
	Slotinr = ChatConfig{BotToken: SlotinrBotToken, ChatID: SlotinrChatID}
	RR9     = ChatConfig{BotToken: RR9BotToken, ChatID: RR9ChatID}
	A777win = ChatConfig{BotToken: A777winBotToken, ChatID: A777winChatID}
	Luckbet = ChatConfig{BotToken: LuckbetBotToken, ChatID: LuckbetChatID}
	Yy6     = ChatConfig{BotToken: Yy6BotToken, ChatID: Yy6ChatID}
	A01Game = ChatConfig{BotToken: A01GameBotToken, ChatID: A01GameChatID}
	P77     = ChatConfig{BotToken: P77BotToken, ChatID: P77ChatID}
	Wininr  = ChatConfig{BotToken: WininrBotToken, ChatID: WininrChatID}
	Ww5     = ChatConfig{BotToken: Ww5BotToken, ChatID: Ww5ChatID}
	A77bet  = ChatConfig{BotToken: A77betBotToken, ChatID: A77betChatID}
	Dd1     = ChatConfig{BotToken: Dd1BotToken, ChatID: Dd1ChatID}
	A5222   = ChatConfig{BotToken: A5222BotToken, ChatID: A5222ChatID}
)

type ChatConfig struct {
	BotToken string
	ChatID   string
}

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

// Alert 定义 Prometheus Alert 数据结构
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// logger 全局日志对象
var logger *zap.Logger

// 初始化日志
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
		BotToken: getEnv("TELEGRAM_BOT_TOKEN", "8108753238:AAGIFeUhr4lWQiYbh3QvNzoZM6x1zfAFgHc"), // 从环境变量中获取 Token
		ChatID:   getEnv("TELEGRAM_CHAT_ID", "-1002428509416"),                                   // 从环境变量中获取 ChatID
		Port:     getEnv("PORT", ":9910"),                                                        // 默认监听 9910 端口
	}

	if config.BotToken == "" || config.ChatID == "" {
		return nil, fmt.Errorf("BotToken 或 ChatID 未设置")
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

			fmt.Println("原始告警如下：", alerts.Alerts)
			fmt.Println("当前告警如下：", message)

			// 同步FP发送
			SendAlertFPMessage(message)

			// 开始发送
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
	status := "🚨 触发"
	if alert.Status == "resolved" {
		status = "✅ 恢复"
	}

	return fmt.Sprintf(
		"*%s告警 - %s*\n*告警名称*: %s\n*告警对象*: %s\n*严重级别*: %s\n*摘要*: %s\n*描述*: %s\n*触发时间*: %s\n*恢复时间*: %s\n*详情链接*: [查看详情](%s)",
		status,
		alert.Labels["environment"],      // 提取环境标签
		alert.Labels["alertname"],        // 提取告警名称
		alert.Labels["host"],             // 提取主机标签
		alert.Labels["severity"],         // 提取严重级别
		alert.Annotations["summary"],     // 提取摘要
		alert.Annotations["description"], // 提取描述
		alert.StartsAt,                   // 告警触发时间
		alert.EndsAt,                     // 告警恢复时间
		// alert.GeneratorURL,               // 告警详情链接
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

// main 入口函数
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

// 分盘关键字列表
var keywords = []string{
	"jeet",
	"p77",
	"yk777",
	"u77",
	"asgame",
	"ez777",
	"slotinr",
	"rr9",
	"a777win",
	"luckbet",
	"yy6",
	"a01game",
	"wininr",
	"ww5",
	"a77bet",
	"dd1",
	"a5222",
}

// 标签区分，发送到对应分盘告警群
func AlertGroupCheckList(message string) (bool, string) {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(message), strings.ToLower(keyword)) {
			return true, keyword
		}
	}
	return false, ""
}

func SendAlertFPMessage(msg string) {
	// 执行检测标签，返回项目名称
	ok, projectName := AlertGroupCheckList(msg)
	if !ok {
		logger.Error("未匹配到对应分盘")
		return
	}

	// 根据项目名称获取对应的 ChatConfig
	var chatConfig ChatConfig
	switch projectName {
	case "jeet":
		chatConfig = Jeet
	case "yk777":
		chatConfig = Yk777
	case "u77":
		chatConfig = U77
	case "asgame":
		chatConfig = Asgame
	case "ez777":
		chatConfig = Ez777
	case "slotinr":
		chatConfig = Slotinr
	case "rr9":
		chatConfig = RR9
	case "a777win":
		chatConfig = A777win
	case "luckbet":
		chatConfig = Luckbet
	case "yy6":
		chatConfig = Yy6
	case "a01game":
		chatConfig = A01Game
	case "p77":
		chatConfig = P77
	case "wininr":
		chatConfig = Wininr
	case "ww5":
		chatConfig = Ww5
	case "a77bet":
		chatConfig = A77bet
	case "dd1":
		chatConfig = Dd1
	case "a5222":
		chatConfig = A5222
	default:
		logger.Error("未找到对应项目名称的 ChatConfig")
		return
	}

	// 发送消息到对应的 Telegram 群组
	err := sendTelegramMessage(chatConfig.BotToken, chatConfig.ChatID, msg)
	if err != nil {
		logger.Error("Failed to send alert to Telegram",
			zap.Error(err),
			zap.String("chat_id", chatConfig.ChatID))
	}
}
