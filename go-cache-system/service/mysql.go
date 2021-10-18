/**
 * @Author: djh
 * @Description:
 * @File:  mysql
 * @Version: 1.0.0
 * @Date: 2021/9/23 19:50
 */

package service

import (
	"cache-system/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func NewMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Configs.Mysql.User,
		config.Configs.Mysql.Password,
		config.Configs.Mysql.Host,
		config.Configs.Mysql.Port,
		config.Configs.Mysql.Database)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Mysql connect err:%s", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("Mysql db set err:%s", err))
	}
	sqlDB.SetMaxIdleConns(config.Configs.Mysql.MaxIdle)
	sqlDB.SetMaxOpenConns(config.Configs.Mysql.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Configs.Mysql.MaxLife))

}
