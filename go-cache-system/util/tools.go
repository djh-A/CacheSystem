/**
 * @Author: djh
 * @Description:
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/9/24 19:27
 */

package util

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"time"
)

const (
	DATEFORMAT      = "2006-01-02"
	TIMESTAMPFORMAT = "2006-01-02 15:04:05"
	Y               = "2006"
	Yy              = "06"
	M               = "01"
	Mm              = "1"
	D               = "02"
	Dd              = "2"
)

//获取2个日期相差几天
func TimeSubDays(t1, t2 time.Time) int {
	if t1.Location().String() != t2.Location().String() {
		return -1
	}
	hours := t1.Sub(t2).Hours()
	if hours <= 0 {
		return -1
	}
	// sub hours less than 24
	if hours < 24 {
		// may same day
		t1y, t1m, t1d := t1.Date()
		t2y, t2m, t2d := t2.Date()
		isSameDay := (t1y == t2y && t1m == t2m && t1d == t2d)

		if isSameDay {

			return 0
		} else {
			return 1
		}

	} else { // equal or more than 24

		if (hours/24)-float64(int(hours/24)) == 0 { // just 24's times
			return int(hours / 24)
		} else { // more than 24 hours
			return int(hours/24) + 1
		}
	}
}

//获取2个日期相差几天
func TimeStrSubDays(t1, t2 string) (float64, time.Time, time.Time) {

	time1, _ := time.Parse(TIMESTAMPFORMAT, t1)
	time2, _ := time.Parse(TIMESTAMPFORMAT, t2)

	return math.Round(time2.Sub(time1).Hours() / 24), time1, time2
}

//获取2个日期相差几天
func TimeStrSubMonths(t1, t2 time.Time) int {

	return int(t2.Month() - t1.Month())
}

//GZIPEn gzip加密
//str:="Hello 蓝影闪电"
//	strGZIPEn:= GZIPEn(str)
//	fmt.Println(strGZIPEn) //加密
//	strGZIPDe,_:=GZIPDe(strGZIPEn)
//	fmt.Println(string(strGZIPDe)) //解密
func GZIPEn(str string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(str)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

//GZIPDe gzip解密
func GZIPDe(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func ParseCacheKey(s string) string {
	re := regexp.MustCompile(`\s+`)
	str := re.ReplaceAllString(s, "")
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
