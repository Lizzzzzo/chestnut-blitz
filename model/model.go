package model

import "time"

// 活动表
type Activity struct {
	ID           int       // 活动ID
	Name         string    // 活动名称
	Desc         string    // 活动描述
	StartTime    time.Time // 活动开始时间
	EndTime      time.Time // 活动结束时间
	ProductID    int       // 活动关联的商品ID
	ProductStock int       // 活动关联的商品库存
}

type Order struct {
	ID         int // 订单ID
	ActivityID int // 订单关联的活动ID
	UserID     int // 订单关联的用户ID
	Status     int // 订单状态
}
