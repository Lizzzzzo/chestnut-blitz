package handler

import (
	"chestnut-blitz/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SecKillReq struct {
	ActivityID int `json:"activity_id"`
	UserID     int `json:"user_id"`
}

func SecKill(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定传参
		var req SecKillReq
		err := c.ShouldBindJSON(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// 1、查询活动是否存在
		db.First(&model.Activity{}, req.ActivityID)

		// 2、查询是否处于活动时间范围内
		db.Where("start_at <= ? and end_at >= ?", time.Now(), time.Now()).Find(&model.Activity{})

		// 3、查询商品库存是否不为0
		db.Where("product_stock > 0").Find(&model.Activity{})

		// 4、查询用户是否未购买
		db.Where("user_id = ?", req.UserID).Find(&model.Order{})

		// 5、创建订单
		order := &model.Order{
			ActivityID: req.ActivityID,
			UserID:     req.UserID,
			Status:     0,
		}
		db.Create(order)

		// 6、等待用户支付
		db.Model(&model.Order{}).Where("user_id = ?", req.UserID).Update("status", 1)
	}
}
