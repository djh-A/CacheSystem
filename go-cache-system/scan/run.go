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
	"strings"
	"sync"
	"time"
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
	EffectiveDate      time.Time
	UpdatedAt          time.Time
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

var wg sync.WaitGroup

func before() {
	service.NewImpala()
	service.NewMysql()

}

func after() {
	service.CloseImpala()

}

func Run() error {
	defer func() {
		if err := recover(); nil != err {
			service.SendSMTPMail(fmt.Sprintf("Run error :%s", err))
			panic(err)
		}
	}()
	before()
	redisPool = service.NewRedisPool()
	eliminationStrategy()
	cache := []CacheSystemConfig{}
	service.DB.Where("isCache = ?", CACHE_IN).Where("heat > ?", 1).Order("heat desc").Limit(config.Configs.Select.Limit).Find(&cache)
	//Distribute the data evenly to each coroutine

	totalLen := len(cache)
	if totalLen == 0 {
		return nil
	}
	goroutine := config.Configs.Select.Goroutine
	if totalLen <= goroutine {
		goroutine = totalLen
	}
	goLen := func() int {
		if totalLen <= goroutine {
			return 1
		}
		return int(math.Ceil(float64(totalLen) / float64(goroutine)))
	}()
	//If the data is small, reset the custom number of coroutines

	signalChan := make(chan byte, goroutine)

	//var goroutine = make(chan int)
	for i := 1; i <= goroutine; i++ {
		sliceStart := (i - 1) * goLen
		sliceEnd := func() int {
			if goLen+sliceStart > totalLen {
				return totalLen
			}
			return goLen + sliceStart

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

func eliminationStrategy() {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	key := "clear:heatClear"
	response, err := redisConn.Do(service.REDIS_GET, key)
	if err != nil {
		Log.Infof("Redis set [%s] err = ", key, err)
	}
	if response == nil {
		heatClear()
		_, _ = redisConn.Do(service.REDIS_SET, key, true)
		_, _ = redisConn.Do(service.REDIS_EXPIRE, key, 43200)
	}
}

//Heat elimination strategy
func heatClear() {
	wg.Add(1)
	go func() {
		service.DB.Where("heat = ?", 1).Delete(CacheSystemConfig{})
		wg.Done()
		Log.Info("Cleaned up!")
	}()
	wg.Wait()
}

func goRun(arr []CacheSystemConfig, signalChan chan<- byte) {
	defer func() {
		if err := recover(); nil != err {
			Log.Error(err)
			signalChan <- 1
			service.SendSMTPMail(fmt.Sprintf("Go run error :%s", err))
		}
	}()

	for _, cache := range arr {

		//If the cache validity period is exceeded, then refresh the cache
		if int(time.Now().Sub(cache.UpdatedAt).Seconds()) >= cache.ExpTime-int(cache.QueryTimeConsuming)-int((time.Duration(config.Configs.Select.Sleep)*time.Second).Seconds()) {
			Log.Infof("ID: %d sql cache expired and will be refreshed", cache.Id)
			if cache.EffectiveDate.Format(util.DATEFORMAT) != time.Now().Format(util.DATEFORMAT) {
				Log.Infof("ID：%d sql cache validity period has expired and will be replace operation", cache.Id)
				//cache.Year()
				//cache.Month()
				cache.Day()
			}

			var result []byte
			result, cache.QueryTimeConsuming = func(t time.Time) ([]byte, float64) {
				//result, _ := service.FetchWithoutKey(cache.Sql)
				result, err := service.QueryAll(cache.Sql)
				if err != nil {
					panic(err)
				}
				//result := make([]interface{}, 7)
				resultJson, err := json.Marshal(result)
				if err != nil {
					return nil, 0.0
				}
				return resultJson, time.Since(t).Seconds()
			}(time.Now())
			redisConn := redisPool.Get()
			defer redisConn.Close()
			_, err := redisConn.Do(service.REDIS_SET, cache.CacheKey, result)
			if err != nil {
				Log.Infof("Redis set err = ", err)
				panic(fmt.Sprintf("Redis set err:%s", err))
			}
			_, err = redisConn.Do(service.REDIS_EXPIRE, cache.CacheKey, cache.ExpTime)
			if err != nil {
				Log.Infof("Redis set expire err = ", err)
				panic(fmt.Sprintf("Redis set expire err:%s", err))
			}
			Log.Info("Redis cache successfully")
			cache.UpdatedAt, _ = time.Parse(util.TIMESTAMPFORMAT, time.Now().Format(util.TIMESTAMPFORMAT))
			cache.EffectiveDate, _ = time.Parse(util.DATEFORMAT, time.Now().Format(util.DATEFORMAT))
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
	Log.Infof("Before Sql: %s", this.Sql)
	regDateField := strings.Join(dateField, "|")
	re := regexp.MustCompile(fmt.Sprintf(`(%s)\s+BETWEEN\s+'(.*)'\s+AND\s+'(.*)'`, regDateField))
	result := re.FindAllStringSubmatch(this.Sql, -1)
	dateMap := make(map[int]map[string]time.Time)
	for i, res := range result {
		if len(res) > 0 {
			startDate, _ := time.Parse(util.TIMESTAMPFORMAT, res[2])
			endDate, _ := time.Parse(util.TIMESTAMPFORMAT, res[3])
			newStartDate := startDate.AddDate(0, 0, 1)
			newEndDate := endDate.AddDate(0, 0, 1)
			date := make(map[string]time.Time)
			date["startDate"] = newStartDate
			date["endDate"] = newEndDate
			dateMap[i] = date
			//替换日期
			this.Sql = strings.Replace(this.Sql, res[0], fmt.Sprintf("%s BETWEEN '%s' AND '%s'",
				res[1],
				newStartDate.Format(util.TIMESTAMPFORMAT),
				newEndDate.Format(util.TIMESTAMPFORMAT)), -1)
		}
	}
	pre := regexp.MustCompile(`AND\s*\(\s*(.*) \)`)
	results := pre.FindAllStringSubmatch(this.Sql, -1)
	for i, res := range results {
		tableAlias := util.GetTableAlias(res[1])
		if len(res) > 0 {
			partition := util.GetPartition(res[1])
			condition := util.SetPartitionWhereCondition(dateMap[i]["startDate"], dateMap[i]["endDate"], partition, tableAlias)
			if partition != "" {
				this.Sql = strings.Replace(this.Sql, res[1], condition, -1)
			}

		}
	}
	this.CacheKey = util.ParseCacheKey(this.Sql)
	Log.Infof("After Sql: %s", this.Sql)
}
