package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response 定义返回的响应结构
type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// ResponseBody 定义响应体的结构
type ResponseBody struct {
	Message     string          `json:"message"`
	Timestamp   string          `json:"timestamp"`
	Environment EnvironmentInfo `json:"environment"`
	Event       interface{}     `json:"event"`
}

// EnvironmentInfo 定义环境变量信息结构
type EnvironmentInfo struct {
	ENV        string            `json:"ENV"`
	AllEnvVars map[string]string `json:"allEnvVars"`
}

// Handler 是Lambda函数的入口点
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	// 打印事件对象到CloudWatch日志
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("事件对象: %s\n", string(eventJSON))

	// 打印环境变量
	fmt.Printf("ENV环境变量: %s\n", os.Getenv("ENV"))

	// 获取所有环境变量
	allEnvVars := make(map[string]string)
	for _, env := range os.Environ() {
		// 解析环境变量，格式为 "KEY=VALUE"
		if len(env) > 0 {
			for i, c := range env {
				if c == '=' {
					key := env[:i]
					value := env[i+1:]
					allEnvVars[key] = value
					break
				}
			}
		}
	}

	// 打印所有环境变量到日志
	allEnvJSON, _ := json.MarshalIndent(allEnvVars, "", "  ")
	fmt.Printf("环境变量: %s\n", string(allEnvJSON))

	// 构建响应体
	responseBody := ResponseBody{
		Message:   "Hello, World! Go Lambda函数已成功执行了。",
		Timestamp: time.Now().Format(time.RFC3339),
		Environment: EnvironmentInfo{
			ENV:        os.Getenv("ENV"),
			AllEnvVars: allEnvVars,
		},
		Event: event,
	}

	// 将响应体转换为JSON字符串
	bodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       `{"error":"Failed to marshal response"}`,
		}, err
	}

	return Response{
		StatusCode: 200,
		Body:       string(bodyJSON),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
