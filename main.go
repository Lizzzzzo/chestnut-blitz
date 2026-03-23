package main

import (
	"chestnut-blitz/model"

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
}
