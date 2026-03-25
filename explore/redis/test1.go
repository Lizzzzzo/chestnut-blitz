package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rdb.Set(ctx, "test:hello", "world", 60).Err()
	if err != nil {
		fmt.Println("创建失败：", err)
		return
	}

	val, err := rdb.Get(ctx, "test:hello").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key 不存在!")
		} else {
			fmt.Println("获取失败：", err)
		}
		return
	} else {
		fmt.Println(val)
	}

	err = rdb.Del(ctx, "test:hello").Err()
	if err != nil {
		fmt.Println("删除失败：", err)
		return
	}
}
