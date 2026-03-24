package model

import (
	"time"

	"gorm.io/gorm"
)

// 活动表
type Activity struct {
	// ID           int       // 活动ID
	gorm.Model
	Name         string    // 活动名称
	Desc         string    // 活动描述
	StartTime    time.Time // 活动开始时间
	EndTime      time.Time // 活动结束时间
	ProductID    int       // 活动关联的商品ID
	ProductStock int       // 活动关联的商品库存
}

type Order struct {
	// ID         int // 订单ID
	gorm.Model
	ActivityID int // 订单关联的活动ID
	UserID     int // 订单关联的用户ID
	ProductID  int // 订单关联的商品ID
	Status     int // 订单状态 0代表待支付 1代表支付成功 2代表支付失败
}
