package main

import (
	"chestnut-blitz/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("chestnut-blittz.db"), &gorm.Config{})
	if err != nil {
		panic("failer to connect database")
	}
	db.AutoMigrate(
		&model.Activity{},
		&model.Order{},
	)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping success",
		})
	})

	r.Run()
}
