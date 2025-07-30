package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// NotificationRequest 通知请求结构 - 通用接口设计
type NotificationRequest struct {
	// 基础信息
	MessageType string `json:"message_type"` // "error", "warning", "info", "success", "alert"
	Title       string `json:"title"`        // 消息标题
	Content     string `json:"content"`      // 消息内容

	// 来源信息
	Source struct {
		ServiceName string `json:"service_name"` // 服务名称
		ModuleName  string `json:"module_name"`  // 模块名称
		Environment string `json:"environment"`  // 环境：dev, test, staging, prod
		Region      string `json:"region"`       // 区域信息
		Version     string `json:"version"`      // 版本信息
	} `json:"source"`

	// 详细信息
	Details struct {
		Level      string                 `json:"level"`       // 级别：DEBUG, INFO, WARN, ERROR, FATAL
		Timestamp  string                 `json:"timestamp"`   // 时间戳，ISO8601格式
		TraceID    string                 `json:"trace_id"`    // 追踪ID
		UserID     string                 `json:"user_id"`     // 用户ID（如果适用）
		RequestID  string                 `json:"request_id"`  // 请求ID
		ErrorCode  string                 `json:"error_code"`  // 错误代码
		StackTrace string                 `json:"stack_trace"` // 堆栈信息
		Metadata   map[string]interface{} `json:"metadata"`    // 额外元数据
	} `json:"details"`

	// 目标配置
	Target struct {
		ChatIDs     []string `json:"chat_ids"`     // 目标群聊ID列表
		ProjectName string   `json:"project_name"` // 项目名称（用于映射群聊）
		Priority    string   `json:"priority"`     // 优先级：low, normal, high, urgent
		AtAll       bool     `json:"at_all"`       // 是否@所有人
		AtUsers     []string `json:"at_users"`     // @指定用户ID列表
	} `json:"target"`

	// 扩展字段（为未来功能预留）
	Extensions struct {
		RetryCount       int                    `json:"retry_count"`       // 重试次数
		CallbackURL      string                 `json:"callback_url"`      // 回调URL
		Tags             []string               `json:"tags"`              // 标签
		CustomFields     map[string]interface{} `json:"custom_fields"`     // 自定义字段
		RateLimitKey     string                 `json:"rate_limit_key"`    // 限流key
		DeduplicationKey string                 `json:"deduplication_key"` // 去重key
	} `json:"extensions"`
}

// NotificationResponse 响应结构
type NotificationResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Details   struct {
		ProcessedAt    string   `json:"processed_at"`
		SentToChats    []string `json:"sent_to_chats"`
		FailedChats    []string `json:"failed_chats"`
		RetryScheduled bool     `json:"retry_scheduled"`
	} `json:"details"`
}

// FeishuClient 飞书客户端
type FeishuClient struct {
	AppID       string
	AppSecret   string
	BaseURL     string
	AccessToken string
	TokenExpire time.Time
	httpClient  *http.Client
}

// FeishuAccessTokenResponse 飞书访问令牌响应
type FeishuAccessTokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	AppAccessToken    string `json:"app_access_token"`
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// FeishuAPIResponse 飞书API通用响应
type FeishuAPIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ChatInfo 群聊信息
type ChatInfo struct {
	ChatID      string `json:"chat_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	OwnerID     string `json:"owner_id"`
	Members     int    `json:"members"`
}

// ChatListResponse 群聊列表响应
type ChatListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items     []ChatInfo `json:"items"`
		PageToken string     `json:"page_token"`
		HasMore   bool       `json:"has_more"`
	} `json:"data"`
}

// NewFeishuClient 创建飞书客户端
func NewFeishuClient() *FeishuClient {
	return &FeishuClient{
		AppID:     os.Getenv("FEISHU_APP_ID"),
		AppSecret: os.Getenv("FEISHU_APP_SECRET"),
		BaseURL:   getEnvWithDefault("FEISHU_BASE_URL", "https://open.feishu.cn"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getAccessToken 获取访问令牌
func (c *FeishuClient) getAccessToken() error {
	if time.Now().Before(c.TokenExpire) && c.AccessToken != "" {
		return nil
	}

	url := fmt.Sprintf("%s/open-apis/auth/v3/tenant_access_token/internal", c.BaseURL)

	payload := map[string]string{
		"app_id":     c.AppID,
		"app_secret": c.AppSecret,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp FeishuAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if tokenResp.Code != 0 {
		return fmt.Errorf("failed to get access token: %s", tokenResp.Msg)
	}

	c.AccessToken = tokenResp.TenantAccessToken
	c.TokenExpire = time.Now().Add(time.Duration(tokenResp.Expire-300) * time.Second) // 提前5分钟过期

	return nil
}

// GetChatList 获取机器人加入的群聊列表
func (c *FeishuClient) GetChatList() ([]ChatInfo, error) {
	if err := c.getAccessToken(); err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/open-apis/im/v1/chats", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page_size", "50")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("API响应: %s\n", string(body))

	var chatListResp ChatListResponse
	if err := json.Unmarshal(body, &chatListResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatListResp.Code != 0 {
		return nil, fmt.Errorf("failed to get chat list: %s", chatListResp.Msg)
	}

	return chatListResp.Data.Items, nil
}

// SendMessage 发送消息到群聊
func (c *FeishuClient) SendMessage(chatID, message string) error {
	if err := c.getAccessToken(); err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/open-apis/im/v1/messages?receive_id_type=chat_id", c.BaseURL)

	contentJSON, err := json.Marshal(map[string]string{
		"text": message,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	requestBody := map[string]interface{}{
		"receive_id": chatID,
		"msg_type":   "text",
		"content":    string(contentJSON),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp FeishuAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResp.Code != 0 {
		return fmt.Errorf("failed to send message: %s", apiResp.Msg)
	}

	return nil
}

// FormatMessage 格式化消息内容
func FormatMessage(req NotificationRequest) string {
	var message strings.Builder

	// 根据消息类型添加图标
	switch req.MessageType {
	case "error":
		message.WriteString("🚨 错误报警\n")
	case "warning":
		message.WriteString("⚠️ 警告通知\n")
	case "info":
		message.WriteString("ℹ️ 信息通知\n")
	case "success":
		message.WriteString("✅ 成功通知\n")
	case "alert":
		message.WriteString("📢 重要告警\n")
	default:
		message.WriteString("📝 系统通知\n")
	}

	// 基础信息
	if req.Title != "" {
		message.WriteString(fmt.Sprintf("标题: %s\n", req.Title))
	}

	// 来源信息
	if req.Source.ServiceName != "" {
		message.WriteString(fmt.Sprintf("服务: %s", req.Source.ServiceName))
		if req.Source.ModuleName != "" {
			message.WriteString(fmt.Sprintf("/%s", req.Source.ModuleName))
		}
		message.WriteString("\n")
	}

	if req.Source.Environment != "" {
		message.WriteString(fmt.Sprintf("环境: %s", req.Source.Environment))
		if req.Source.Region != "" {
			message.WriteString(fmt.Sprintf("(%s)", req.Source.Region))
		}
		message.WriteString("\n")
	}

	// 详细信息
	if req.Details.Level != "" {
		message.WriteString(fmt.Sprintf("级别: %s\n", req.Details.Level))
	}

	if req.Details.Timestamp != "" {
		message.WriteString(fmt.Sprintf("时间: %s\n", req.Details.Timestamp))
	} else {
		message.WriteString(fmt.Sprintf("时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	}

	// 内容
	if req.Content != "" {
		message.WriteString(fmt.Sprintf("内容: %s\n", req.Content))
	}

	// 追踪信息
	if req.Details.TraceID != "" {
		message.WriteString(fmt.Sprintf("追踪ID: %s\n", req.Details.TraceID))
	}

	if req.Details.RequestID != "" {
		message.WriteString(fmt.Sprintf("请求ID: %s\n", req.Details.RequestID))
	}

	if req.Details.ErrorCode != "" {
		message.WriteString(fmt.Sprintf("错误码: %s\n", req.Details.ErrorCode))
	}

	// 堆栈信息
	if req.Details.StackTrace != "" {
		message.WriteString(fmt.Sprintf("堆栈信息:\n%s\n", req.Details.StackTrace))
	}

	// 版本信息
	if req.Source.Version != "" {
		message.WriteString(fmt.Sprintf("版本: %s\n", req.Source.Version))
	}

	return message.String()
}

// getChatIDsByProject 根据项目名获取群聊ID列表
func getChatIDsByProject(projectName string) []string {
	// 从环境变量中读取项目到群聊的映射配置
	// 格式: PROJECT_CHAT_MAPPING={"project1":"chat_id1","project2":"chat_id2,chat_id3"}
	mappingStr := os.Getenv("PROJECT_CHAT_MAPPING")
	if mappingStr == "" {
		// 返回默认群聊ID
		defaultChatID := os.Getenv("DEFAULT_CHAT_ID")
		if defaultChatID != "" {
			return []string{defaultChatID}
		}
		return []string{}
	}

	var mapping map[string]string
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		fmt.Printf("Error parsing PROJECT_CHAT_MAPPING: %v\n", err)
		return []string{}
	}

	if chatIDs, exists := mapping[projectName]; exists {
		// 支持逗号分隔的多个群聊ID
		return strings.Split(chatIDs, ",")
	}

	// 如果没有找到对应项目，返回默认群聊
	if defaultChatIDs, exists := mapping["default"]; exists {
		return strings.Split(defaultChatIDs, ",")
	}

	return []string{}
}

// Handler Lambda函数处理器
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("收到请求: %s %s\n", request.HTTPMethod, request.Path)
	fmt.Printf("请求头: %+v\n", request.Headers)
	fmt.Printf("请求体: %s\n", request.Body)

	// 健康检查
	if request.HTTPMethod == "GET" && request.Path == "/health" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"status":"healthy","service":"feishu-notification-service","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`,
		}, nil
	}

	// 获取群聊列表
	if request.HTTPMethod == "GET" && request.Path == "/chats" {
		return handleGetChats()
	}

	// 只接受POST请求
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"Method not allowed"}`,
		}, nil
	}

	// 解析请求体
	var notificationReq NotificationRequest
	if err := json.Unmarshal([]byte(request.Body), &notificationReq); err != nil {
		fmt.Printf("解析请求失败: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"Invalid JSON format"}`,
		}, nil
	}

	// 验证必需字段
	if notificationReq.Content == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"Content is required"}`,
		}, nil
	}

	// 创建飞书客户端
	feishuClient := NewFeishuClient()
	if feishuClient.AppID == "" || feishuClient.AppSecret == "" {
		fmt.Println("飞书配置未设置")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"Feishu configuration not set"}`,
		}, nil
	}

	// 确定目标群聊ID
	var targetChatIDs []string
	if len(notificationReq.Target.ChatIDs) > 0 {
		targetChatIDs = notificationReq.Target.ChatIDs
	} else if notificationReq.Target.ProjectName != "" {
		targetChatIDs = getChatIDsByProject(notificationReq.Target.ProjectName)
	} else {
		// 使用默认群聊
		defaultChatID := os.Getenv("DEFAULT_CHAT_ID")
		if defaultChatID != "" {
			targetChatIDs = []string{defaultChatID}
		}
	}

	if len(targetChatIDs) == 0 {
		fmt.Println("未找到目标群聊ID")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"No target chat IDs found"}`,
		}, nil
	}

	// 格式化消息
	formattedMessage := FormatMessage(notificationReq)
	fmt.Printf("格式化消息: %s\n", formattedMessage)

	// 发送消息到所有目标群聊
	var sentToChats []string
	var failedChats []string

	for _, chatID := range targetChatIDs {
		fmt.Printf("发送消息到群聊: %s\n", chatID)
		if err := feishuClient.SendMessage(strings.TrimSpace(chatID), formattedMessage); err != nil {
			fmt.Printf("发送消息失败 (群聊: %s): %v\n", chatID, err)
			failedChats = append(failedChats, chatID)
		} else {
			fmt.Printf("发送消息成功 (群聊: %s)\n", chatID)
			sentToChats = append(sentToChats, chatID)
		}
	}

	// 构建响应
	response := NotificationResponse{
		Success:   len(sentToChats) > 0,
		MessageID: fmt.Sprintf("%d", time.Now().UnixNano()),
		Message:   fmt.Sprintf("Sent to %d chats, failed %d chats", len(sentToChats), len(failedChats)),
		Timestamp: time.Now().Format(time.RFC3339),
		Details: struct {
			ProcessedAt    string   `json:"processed_at"`
			SentToChats    []string `json:"sent_to_chats"`
			FailedChats    []string `json:"failed_chats"`
			RetryScheduled bool     `json:"retry_scheduled"`
		}{
			ProcessedAt:    time.Now().Format(time.RFC3339),
			SentToChats:    sentToChats,
			FailedChats:    failedChats,
			RetryScheduled: false,
		},
	}

	responseBody, _ := json.Marshal(response)

	statusCode := 200
	if len(sentToChats) == 0 {
		statusCode = 500
	} else if len(failedChats) > 0 {
		statusCode = 207 // 部分成功
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

// handleGetChats 处理获取群聊列表请求
func handleGetChats() (events.APIGatewayProxyResponse, error) {
	fmt.Println("处理获取群聊列表请求")

	// 创建飞书客户端
	feishuClient := NewFeishuClient()
	if feishuClient.AppID == "" || feishuClient.AppSecret == "" {
		fmt.Println("飞书配置未设置")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"success":false,"message":"Feishu configuration not set"}`,
		}, nil
	}

	// 获取群聊列表
	chats, err := feishuClient.GetChatList()
	if err != nil {
		fmt.Printf("获取群聊列表失败: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf(`{"success":false,"message":"Failed to get chat list: %s"}`, err.Error()),
		}, nil
	}

	// 构建响应
	response := struct {
		Success   bool       `json:"success"`
		Message   string     `json:"message"`
		Timestamp string     `json:"timestamp"`
		Count     int        `json:"count"`
		Chats     []ChatInfo `json:"chats"`
	}{
		Success:   true,
		Message:   fmt.Sprintf("Successfully retrieved %d chats", len(chats)),
		Timestamp: time.Now().Format(time.RFC3339),
		Count:     len(chats),
		Chats:     chats,
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	lambda.Start(Handler)
}
