/**
 * @Author: djh
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2021/9/23 10:52
 */

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type (
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
		MaxIdle  int `json:"max_idle"` //设置连接池中空闲连接的最大数量  
		MaxOpen  int `json:"max_open"` //打开数据库连接的最大数量
		MaxLife  int `json:"max_life"` //连接可复用的最大时间，秒为单位
	}
	Impala struct {
		Host     string
		Port     int
		Database string
	}
	Redis struct {
		Host        string
		Password    string
		Port        int
		MaxIdle     int `json:"max_idle"`     //设置连接池中空闲连接的最大数量
		IdleTimeout int `json:"idle_timeout"` //连接可复用的最大时间，秒为单位
		Database    int
	}
	Select struct {
		Goroutine int
		Limit     int
		Sleep     int
		Env       string
	}
	Email struct {
		Host     string
		Port     int
		User     string
		Password string
		To       []string
		Subject  string
	}
)

type Config struct {
	Mysql  Mysql
	Redis  Redis
	Select Select
	Impala Impala
	Email  Email
}

var Configs *Config

func init() {

}

func NewConfig() *Config {
	var conf string
	flag.StringVar(&conf, "conf", "./config/conf.json", "-conf path")
	flag.Parse()
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		fmt.Println("fail to open the file,", err)
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println("Serialization failed,", err)
	}
	Configs = &config
	fmt.Printf("Config path:%s\n", conf)
	fmt.Println(Configs)
	return Configs
}
