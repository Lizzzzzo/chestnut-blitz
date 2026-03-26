package handler

import (
	"chestnut-blitz/model"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SecKillReq struct {
	ActivityID int `json:"activity_id"`
	UserID     int `json:"user_id"`
}

const luaSecKill = `
local stockKey = KEYS[1]
local userKey = KEYS[2]
local userID = ARGV[1]

if redis.call('SISMEMBER', userKey, userID) == 1 then
	return -1
end

local stock = redis.call('GET', stockKey)
if not stock or tonumber(stock) <= 0 then
	return 0
end

redis.call('DECR', stockKey)
redis.call('SADD', userKey, userID)

return 1
`

const luaReturnStock = `
local stockKey = KEYS[1]
local userKey = KEYS[2]
local userID = ARGV[1]

local stock = redis.call('GET', stockKey)
if not stock then 
	return -1
end

if redis.call('SISMEMBER', userKey, userID) == 1 then
	redis.call('INCR', stockKey)
	redis.call('SREM', userKey, userID)
	return 1
else
	return 0
end
`

func SecKill(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定传参
		var req SecKillReq
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// 1、查询活动是否存在
		var act model.Activity
		err = db.First(&act, req.ActivityID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"err": "活动不存在!"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"err": "活动查询失败!"})
			return
		}

		// 2、查询是否处于活动时间范围内
		if act.StartTime.After(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"err": "活动未开始!"})
			return
		}
		if act.EndTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"err": "活动已结束!"})
			return
		}

		// 3、查询商品库存是否不为0
		// if act.ProductStock <= 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "活动商品已售罄!"})
		// 	return
		// }
		// ctx := context.Background()
		// stockKey := "stock:" + strconv.Itoa(req.ActivityID)
		// stockStr, err := rdb.Get(ctx, stockKey).Result()
		// if err != nil {
		// 	if err == redis.Nil {
		// 		c.JSON(http.StatusBadRequest, gin.H{"err": "商品库存不存在!"})
		// 		return
		// 	}
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "获取商品库存失败!"})
		// 	return
		// }
		// stock, err := strconv.Atoi(stockStr)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "字符串转换失败!"})
		// 	return
		// }
		// fmt.Printf("目前商品总库存为: %d\n", stock)
		// if stock <= 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "活动商品已售罄!"})
		// 	return
		// }

		// 4、查询用户是否重复购买
		// var order model.Order
		// err = db.Where("user_id = ? and activity_id = ?", req.UserID, req.ActivityID).
		// 	First(&order).Error
		// if err == nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "您已参加过该活动!"})
		// 	return
		// }
		// if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "查询用户订单失败!"})
		// 	return
		// }

		// 5、商品库存-1
		// result := db.Model(&act).
		// 	Where("product_stock > 0").
		// 	Update("product_stock", gorm.Expr("product_stock - ?", 1))
		// if result.Error != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "商品库存扣减失败!"})
		// 	return
		// }
		// if result.RowsAffected == 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "商品库存不足!"})
		// 	return
		// }

		ctx := context.Background()
		stockKey := "stock:" + strconv.Itoa(req.ActivityID)
		userKey := "user:" + strconv.Itoa(req.ActivityID)
		luaRes, err := rdb.Eval(ctx, luaSecKill, []string{stockKey, userKey}, req.UserID).Int()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "lua脚本获取失败!"})
			return
		}
		fmt.Println("Lua 返回：", luaRes)
		if luaRes == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"err": "您已参加过该活动!"})
			return
		}
		if luaRes == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"err": "活动商品已售罄!"})
			return
		}

		// 6、创建订单
		// order := model.Order{
		// 	ActivityID: req.ActivityID,
		// 	UserID:     req.UserID,
		// 	ProductID:  act.ProductID,
		// 	Status:     0,
		// }
		// err = db.Create(&order).Error
		// if err != nil {
		// 	// 回滚库存
		// 	db.Model(&act).Update("product_stock", gorm.Expr("product_stock + ?", 1))
		// 	c.JSON(http.StatusBadRequest, gin.H{"err": "订单创建失败!"})
		// 	return
		// }
		go func() {
			order := model.Order{
				ActivityID: req.ActivityID,
				UserID:     req.UserID,
				ProductID:  act.ProductID,
				Status:     0,
			}
			err = db.Create(&order).Error
			if err != nil {
				// 回滚库存
				luaRes, err = rdb.Eval(ctx, luaReturnStock, []string{stockKey, userKey}, req.UserID).Int()
				if err != nil {
					fmt.Println("订单创建失败!")
					return
				}
				switch luaRes {
				case -1:
					fmt.Printf("商品库存键不存在! 活动ID: %d\n", req.ActivityID)
					return
				case 0:
					fmt.Printf("用户%d不在集合中，无需回滚\n", req.UserID)
				case 1:
					fmt.Printf("用户%d订单创建失败, 已回滚库存\n", req.UserID)
				}
				return
			}
			fmt.Printf("用户%d订单创建成功!", req.UserID)
		}()

		// c.JSON(http.StatusOK, gin.H{
		// 	"msg":      "下单成功",
		// 	"order_id": order.ID,
		// })

		c.JSON(http.StatusOK, gin.H{
			"msg": "排队中",
		})
	}
}
