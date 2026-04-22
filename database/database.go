// Package database @Author: youngalone [2023/8/4]
package database

import (
	"bytedancedemo/database/mysql"
	"bytedancedemo/database/redis"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes both MySQL and Redis connections
func Init() {
	mysql.Init()
	redis.Init()
	DB = mysql.DB
}
