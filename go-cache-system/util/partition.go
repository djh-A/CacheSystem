/**
 * @Author: djh
 * @Description:
 * @File:  partition
 * @Version: 1.0.0
 * @Date: 2021/10/25 11:55
 */

package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	YyMmDd = "YMD"
	YyMm   = "YM"
	Yyy    = "Y"
)

func SetPartitionWhereCondition(startDate, endDate time.Time, partition, tableAlias string) string {
	switch partition {
	case YyMmDd:
		return yearMonthDayPartitionCondition(startDate, endDate, tableAlias)
	case YyMm:
		yearMonthPartitionCondition(startDate, endDate, tableAlias)
	case Yyy:
		yearPartitionCondition(startDate, endDate, tableAlias)
	}
	return ""
}

func yearMonthDayPartitionCondition(startDate, endDate time.Time, tableAlias string) string {
	ys, _ := strconv.Atoi(startDate.Format(Yy))
	ye, _ := strconv.Atoi(endDate.Format(Yy))
	ms, _ := strconv.Atoi(startDate.Format(M))
	me, _ := strconv.Atoi(endDate.Format(M))
	ds := startDate.Day()
	de := endDate.Day()
	if ys == ye && ms == me {
		return fmt.Sprintf(`%syear = %s AND %smonth = %s AND %sday BETWEEN %s AND %s`,
			tableAlias, startDate.Format(Yy), tableAlias, startDate.Format(M), tableAlias, startDate.Format(D), endDate.Format(D))
	}
	part := make([]string, 0)
	for yy := ys; yy <= ye; yy++ {
		minM := func() int {
			if yy == ys {
				return ms
			}
			return 1
		}()
		maxM := func() int {
			if yy == ye {
				return me
			}
			return 12
		}()
		for ; minM <= maxM; minM++ {
			minD := func() int {
				if yy == ys && minM == ms {
					return ds
				}
				return 1
			}()
			maxD := func() int {
				if yy == ye && minM == me {
					return de
				}
				return 31
			}()
			part = append(part, fmt.Sprintf("(%syear = %02d AND %smonth = %02d AND %sday BETWEEN %02d AND %02d)",
				tableAlias, yy, tableAlias, minM, tableAlias, minD, maxD))
		}
	}
	return strings.Join(part, " OR ")
}

func yearMonthPartitionCondition(startDate, endDate time.Time, tableAlias string) string {
	if startDate.Year() == endDate.Year() {
		return fmt.Sprintf("%syear = %s AND %smonth BETWEEN %s AND %s",
			tableAlias, startDate.Format(Yy), tableAlias, startDate.Format(M), endDate.Format(M))
	}

	ys, _ := strconv.Atoi(startDate.Format(Yy))
	ye, _ := strconv.Atoi(endDate.Format(Yy))
	sql := fmt.Sprintf("( (%syear = %s AND %smonth BETWEEN %s AND 12)",
		tableAlias, ys, tableAlias, startDate.Format(M))
	for i := ys + 1; i < ye; i++ {
		sql += fmt.Sprintf("OR (%syear = %d AND %smonth BETWEEN 1 AND 12)", tableAlias, i, tableAlias)
	}
	sql += fmt.Sprintf("OR (%syear = %s AND %smonth BETWEEN 1 AND %s))", tableAlias, ye, tableAlias, endDate.Format(M))

	return sql
}

func yearPartitionCondition(startDate, endDate time.Time, tableAlias string) string {
	if startDate.Year() == endDate.Year() {
		return fmt.Sprintf("%syear = %s", tableAlias, startDate.Format(Yy))
	}
	return fmt.Sprintf("%syear BETWEEN %s AND %s", tableAlias, startDate.Format(Yy), endDate.Format(Yy))

}

func GetPartition(sql string) string {
	return year(sql) + month(sql) + day(sql)
}

func GetTableAlias(sql string) string {
	re := regexp.MustCompile(`(\w+)\.year`)
	findString := re.FindStringSubmatch(sql)
	if len(findString) > 0 {
		return fmt.Sprintf("%s.", findString[1])
	}
	return ""
}

//Year
func year(sql string) string {
	re := regexp.MustCompile(`year`)
	matchString := re.MatchString(sql)
	if matchString {
		return "Y"
	}
	return ""
}

//Month
func month(sql string) string {
	re := regexp.MustCompile(`month`)
	matchString := re.MatchString(sql)
	if matchString {
		return "M"
	}
	return ""
}

//Day
func day(sql string) string {
	re := regexp.MustCompile(`day`)
	matchString := re.MatchString(sql)
	if matchString {
		return "D"
	}
	return ""
}
