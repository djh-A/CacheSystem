/**
 * @Author: djh
 * @Description:
 * @File:  log
 * @Version: 1.0.0
 * @Date: 2021/10/11 16:28
 */

package util

import (
	"fmt"
	"time"
)

type Log struct {
}

func (this Log) Infof(format string, a ...interface{}) {
	timeInfo := fmt.Sprintf("time=%s", time.Now().String())
	level := "level=info"
	msg := fmt.Sprintf("msg="+format, a)
	f := fmt.Sprintf("%s %s %s", timeInfo, level, msg)
	fmt.Println(f)
}

func (this Log) Info(a ...interface{}) {
	timeInfo := fmt.Sprintf("time=%s", time.Now().String())
	level := "level=info"
	msg := fmt.Sprintf("msg=%s", a)
	f := fmt.Sprintf("%s %s %s", timeInfo, level, msg)
	fmt.Println(f)
}

func (this Log) Error(a ...interface{}) {
	timeInfo := fmt.Sprintf("time=%s", time.Now().String())
	level := "level=error"
	msg := fmt.Sprintf("msg=%s", a)
	f := fmt.Sprintf("%s %s %s", timeInfo, level, msg)
	fmt.Println(f)
}
