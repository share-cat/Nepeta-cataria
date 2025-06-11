package main

import (
	"fmt"
	"log"
	"net/http"
	"share-cat/Nepeta-cataria/pkg/catwheel"
	"time"
)

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:3200/health")
	if err != nil {
		log.Printf("無法連接到服務: %v", err)
		return
	}
	resp.Body.Close()
	log.Printf("服務連接正常，狀態碼: %d", resp.StatusCode)

	catwheelClient := catwheel.NewClient("http://localhost:3200")

	options := catwheel.DefaultExecuteOptions().
		WithLanguage(catwheel.LanguagePython).
		WithTimeout(10).
		WithMemoryLimit(128 * 1024 * 1024)

	log.Printf("執行代碼，選項: %+v\n", options)

	response, err := catwheelClient.ExecuteWithOptions("print('hello world')", options)
	if err != nil {
		log.Fatalf("執行代碼失敗: %v", err)
	}

	fmt.Printf("執行 ID: %s\n", response.ExecutionID)
	fmt.Printf("狀態: %s\n", response.Result.Status)
	fmt.Printf("輸出: %s\n", response.Result.Stdout)
	fmt.Printf("錯誤: %s\n", response.Result.Stderr)

	execution, err := catwheelClient.GetExecution(response.ExecutionID)
	if err != nil {
		log.Printf("獲取執行詳情失敗: %v", err)
	} else {
		fmt.Printf("執行時間: %v\n", execution.Result.Duration())
	}

	executions, err := catwheelClient.ListExecutions()
	if err != nil {
		log.Printf("列出執行記錄失敗: %v", err)
	} else {
		fmt.Printf("總執行次數: %d\n", len(executions.Executions))
	}

	cleanup, err := catwheelClient.CleanupExecution(response.ExecutionID)
	if err != nil {
		log.Printf("清理執行記錄失敗: %v", err)
	} else {
		fmt.Printf("清理訊息: %s\n", cleanup.Message)
	}

	health, err := catwheelClient.HealthCheck()
	if err != nil {
		log.Printf("健康檢查失敗: %v", err)
	} else {
		fmt.Printf("健康狀態: %s\n", health.Status)
	}
}
