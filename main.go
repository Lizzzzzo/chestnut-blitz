package main

import (
	"chestnut-blitz/model"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Liwenhui@tcp(127.0.0.1:3306)/chestnut-blitz?charset=utf8mb4&parserTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s&maxAllowedPacket=0"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failer to connect database")
	}
	db.AutoMigrate(
		&model.Activity{},
		&model.Order{},
	)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping success",
		})
	})

	r.Run()
}
