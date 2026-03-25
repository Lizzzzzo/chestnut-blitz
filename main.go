package main

import (
	"chestnut-blitz/handler"
	"chestnut-blitz/model"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. 连接 MySQL
	dsn := "root:lwh260119@tcp(127.0.0.1:3306)/chestnut_blitz?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s&maxAllowedPacket=0"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect database")
	}
	db.AutoMigrate(
		&model.Activity{},
		&model.Order{},
	)

	// 2. 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis 连接失败：", err)
		return
	}
	fmt.Println("redis 连接成功! 返回：", pong)

	r := gin.Default()
	r.POST("/seckill", handler.SecKill(db))

	r.Run()
}
