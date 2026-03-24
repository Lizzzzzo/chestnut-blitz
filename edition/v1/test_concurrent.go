package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// 并发测试配置
const (
	targetUrl    = "http://localhost:8080/seckill" // 目标请求url
	totalReq     = 100000                          // 总请求数
	concurrency  = 10000                           // 并发数
	activityID   = 1                               // 活动ID
	productStock = 100                             // 商品总库存，用于验证是否超卖
)

type SecKillReq struct {
	ActivityID int `json:"activity_id"`
	UserID     int `json:"user_id"`
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	successCount, failCount := 0, 0

	// 控制并发数
	sem := make(chan struct{}, concurrency)

	for i := 1; i <= totalReq; i++ {
		wg.Add(1)
		sem <- struct{}{} // 获取令牌
		go func(userID int) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			reqBody := SecKillReq{
				ActivityID: activityID,
				UserID:     userID,
			}
			jsonData, _ := json.Marshal(reqBody)

			resp, err := http.Post(targetUrl, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				mu.Lock()
				failCount++
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			} else {
				mu.Lock()
				failCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("成功：%d次，失败：%d次\n", successCount, failCount)
}

// 数值设置：
// 商品库存 100
// 并发数 10000
// 总请求数 100000
//
// 观察结果：
// 成功：80次，失败：99920次
//
// Gin日志：
// 2026/03/24 19:23:30 /Users/liwenhui/WorkSpace/Chestnut/chestnut-blitz/handler/handler.go:30 SLOW SQL >= 200ms
// Error 1040: Too many connections
// dial tcp 127.0.0.1:3306: i/o timeout
//
// 数据库：
// product_stock 为 0
// 订单创建总数 80
