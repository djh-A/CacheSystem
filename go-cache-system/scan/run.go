/**
 * @Author: djh
 * @Description:
 * @File:  run
 * @Version: 1.0.0
 * @Date: 2021/9/23 18:00
 */

package scan

import (
	"cache-system/config"
	"cache-system/service"
	"cache-system/util"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"regexp"
	"strconv"
	. "time"
)

const CACHE_IN = 1

type filter interface {
	Year()
	Month()
	Day()
	Date()
}

type CacheSystemConfig struct {
	Id                 int64
	Sql                string
	CacheKey           string `gorm:"column:cacheKey"`
	IsCache            int    `gorm:"column:isCache"`
	ExpTime            int    `gorm:"column:expTime"`
	DataSize           float32
	CacheTimeConsuming float64
	QueryTimeConsuming float64
	Heat               int64
	EffectiveDate      Time
	UpdatedAt          Time
	Route              string
}

var dateField = []string{
	"the_date",
	"the_date_time",
	"register_time",
	"last_login_time",
	"last_pay_time",
}

var Log = util.Log{}

var redisPool *redis.Pool

func before() {
	service.NewImpala()
	service.NewMysql()

}

func after() {
	service.CloseImpala()

}

func Run() error {
	before()
	redisPool = service.NewRedisPool()
	cache := []CacheSystemConfig{}
	service.DB.Where("isCache", CACHE_IN).Order("heat desc").Limit(config.Configs.Select.Limit).Find(&cache)
	//Distribute the data evenly to each coroutine
	totalLen := len(cache)
	goroutine := config.Configs.Select.Goroutine
	goLen := func() int {
		if totalLen <= goroutine {
			return 1
		}
		return int(math.Ceil(float64(totalLen) / float64(goroutine)))
	}()
	//If the data is small, reset the custom number of coroutines
	if goLen <= 1 {
		goroutine = 1
	}
	signalChan := make(chan byte, goroutine)
	//var goroutine = make(chan int)
	for i := 1; i <= goroutine; i++ {
		sliceStart := (i - 1) * goLen
		sliceEnd := func() int {
			if goLen+sliceStart > totalLen {
				return totalLen
			} else {
				return goLen + sliceStart
			}
		}()

		go goRun(cache[sliceStart:sliceEnd], signalChan)
	}

	i := 1
	for _ = range signalChan {
		if i >= goroutine {
			close(signalChan)
			after()
		}
		i++
	}
	return nil
}

func goRun(arr []CacheSystemConfig, signalChan chan<- byte) {
	defer func() {
		if err := recover(); nil != err {
			Log.Infof("Run error :%s", err)
		}
	}()
	for _, cache := range arr {

		//If the cache validity period is exceeded, then refresh the cache
		if int(Now().Sub(cache.UpdatedAt).Seconds()) >= cache.ExpTime-int(cache.QueryTimeConsuming)-int((Duration(config.Configs.Select.Sleep)*Second).Seconds()) {
			Log.Infof("ID: %d sql cache expired and will be refreshed", cache.Id)
			if cache.EffectiveDate.Format(util.DATEFORMAT) != Now().Format(util.DATEFORMAT) {
				Log.Infof("ID：%d sql cache validity period has expired and will be replace operation", cache.Id)
				//cache.Year()
				//cache.Month()
				cache.Day()
				cache.Date()
			}

			var result []byte
			result, cache.QueryTimeConsuming = func(t Time) ([]byte, float64) {
				//result, _ := service.FetchWithoutKey(cache.Sql)
				result, _ := service.QueryAll(cache.Sql)
				//result := make([]interface{}, 7)
				resultJson, err := json.Marshal(result)
				if err != nil {
					return nil, 0.0
				}
				return resultJson, Since(t).Seconds()
			}(Now())
			redisConn := redisPool.Get()
			defer redisConn.Close()
			_, err := redisConn.Do("Set", cache.CacheKey, result)
			if err != nil {
				Log.Infof("Redis set err = ", err)
				panic(fmt.Sprintf("Redis set err:%s", err))
			}
			_, err = redisConn.Do("expire", cache.CacheKey, cache.ExpTime)
			if err != nil {
				Log.Infof("Redis set expire err = ", err)
				panic(fmt.Sprintf("Redis set expire err:%s", err))
			}
			Log.Info("Redis cache successfully")
			cache.UpdatedAt, _ = Parse(util.TIMESTAMPFORMAT, Now().Format(util.TIMESTAMPFORMAT))
			cache.EffectiveDate, _ = Parse(util.DATEFORMAT, Now().Format(util.DATEFORMAT))
			cache.DataSize = float32(len(result) * 8 / 1024)
			service.DB.Save(&cache)
			Log.Info("Cached data updated successfully")
		}
	}

	signalChan <- 1
}

//处理日期字段
//
func (this *CacheSystemConfig) Date() {
	for _, field := range dateField {
		re := regexp.MustCompile(fmt.Sprintf(`%s\s+BETWEEN\s+'(.*)'\s+AND\s+'(.*)'`, field))
		result := re.FindStringSubmatch(this.Sql)
		if len(result) > 0 {
			_, oldStartDate, oldEndDate := util.TimeStrSubDays(result[1], result[2])
			startDate := oldStartDate.AddDate(0, 0, 1).Format("2006-01-02 00:00:00")
			endDate := oldEndDate.AddDate(0, 0, 1).Format(util.DATEFORMAT)
			this.Sql = re.ReplaceAllStringFunc(this.Sql, func(s string) string {
				return fmt.Sprintf("%s BETWEEN '%s' AND '%s'",
					field,
					startDate,
					fmt.Sprintf("%s 23:59:59", endDate))
			})

		}
	}
}

func (this *CacheSystemConfig) Year() {
	//re := regexp.MustCompile(`year\s+=\s+(\d*)\s*`)
	//result := re.FindStringSubmatch(this.Sql)
}

func (this *CacheSystemConfig) Month() {
	//re := regexp.MustCompile(`month\s+BETWEEN\s+(\d*)\s+AND\s+(\d*)`)
	//result := re.FindStringSubmatch(this.Sql)
	//fmt.Println(result)
	//os.Exit(1)
}

func (this *CacheSystemConfig) Day() {
	re := regexp.MustCompile(`day\s+BETWEEN\s+(\d*)\s+AND\s+(\d*)`)
	result := re.FindStringSubmatch(this.Sql)
	if len(result) > 0 {
		startDay, _ := strconv.Atoi(result[1])
		endDay, _ := strconv.Atoi(result[1])
		this.Sql = re.ReplaceAllStringFunc(this.Sql, func(s string) string {
			return fmt.Sprintf("day BETWEEN %d AND %d", startDay+1, endDay+1)
		})
	}

}
