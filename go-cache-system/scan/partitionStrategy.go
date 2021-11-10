///**
// * @Author: djh
// * @Description:
// * @File:  partitionStrategy
// * @Version: 1.0.0
// * @Date: 2021/10/22 18:05
// */
//
package scan
//
//import (
//	"cache-system/util"
//	"fmt"
//	"os"
//	"regexp"
//	"strconv"
//	"time"
//)
//
//type Strategy interface {
//	ReplacePartition(string, [][]string) string
//}
//
//type Property struct {
//	StartDate time.Time
//	EndDate   time.Time
//	Reg       *regexp.Regexp
//}
//
//type CurrentMonth struct {
//	Property Property
//}
//
//type CrossOneMonth struct {
//	Property Property
//}
//
//type CrossManyMonth struct {
//	Property Property
//}
//
//func NewStrategy(property Property, sql string, result [][]string) string {
//	switch util.TimeStrSubMonths(property.StartDate, property.EndDate) {
//	case 0:
//		month := CurrentMonth{property}
//		return month.ReplacePartition(sql, result)
//	case 1:
//		month := CrossOneMonth{property}
//		return month.ReplacePartition(sql, result)
//	default:
//		month := CrossManyMonth{property}
//		return month.ReplacePartition(sql, result)
//	}
//}
//
//func (this *CurrentMonth) ReplacePartition(sql string, result [][]string) string {
//	return ""
//}
//
////夸一个月，开始日期的开始天数加一天，结束天数不变，结束日期开始天数不变，结束天数加一天
//func (this *CrossOneMonth) ReplacePartition(sql string, result [][]string) string {
//
//	for _, res := range result {
//		year, _ := strconv.Atoi("20" + res[1])
//		month, _ := strconv.Atoi(res[2])
//		var startDay, endDay int
//		if len(res) == 5 {
//			startDay, _ = strconv.Atoi(res[3])
//			endDay, _ = strconv.Atoi(res[4])
//		} else {
//			startDay, _ = strconv.Atoi(res[3])
//			endDay, _ = strconv.Atoi(res[4])
//		}
//		oldStartDate := time.Date(year, time.Month(month), startDay, 0, 0, 0, 0, time.UTC)
//		oldEndDate := time.Date(year, time.Month(month), endDay, 0, 0, 0, 0, time.UTC)
//		newStartDate := oldStartDate.AddDate(0, 0, 1)
//		newEndDate := oldEndDate.AddDate(0, 0, 1)
//		//month, _ := strconv.Atoi(res[2])
//		sql = this.Property.Reg.ReplaceAllStringFunc(sql, func(s string) string {
//			if len(res) == 5 {
//				return fmt.Sprintf("year = %s AND month = %s AND day BETWEEN %s AND %s",
//					this.Property.StartDate.Format(util.Yy),
//					this.Property.StartDate.Format(util.M),
//					this.Property.StartDate.Format(util.D),
//					res[4])
//			}
//			//return fmt.Sprintf("year = %d AND month BETWEEN %d AND %d AND day BETWEEN %d AND %d", newStartDate.Year(), newStartDate.Month(), newEndDate.Month(), newStartDate.Day(), newEndDate.Day())
//			return fmt.Sprintf("year = %s AND month BETWEEN %s AND %s AND day BETWEEN %s AND %s",
//				this.Property.StartDate.Format(util.Yy),
//				this.Property.StartDate.Format(util.M),
//				this.Property.StartDate.Format(util.M),
//				this.Property.StartDate.Format(util.D),
//				this.Property.StartDate.Format(util.D))
//		})
//		fmt.Println(sql, res)
//		os.Exit(1)
//	}
//	//return sql
//	//if int(this.Property.StartDate.Month()) > month {
//	//
//	//}
//	//if len(res) == 5 {
//	//
//	//}
//
//	fmt.Println(sql, result)
//	os.Exit(1)
//	return sql
//}
//
//func (this *CrossManyMonth) ReplacePartition(sql string, result [][]string) string {
//	return ""
//}
