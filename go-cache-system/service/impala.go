/**
 * @Author: djh
 * @Description:
 * @File:  impala
 * @Version: 1.0.0
 * @Date: 2021/10/12 10:41
 */

package service

import (
	"cache-system/config"
	"fmt"
	"github.com/koblas/impalathing"
	"sync"
)

var impala *impalathing.Connection
var mu sync.Mutex

func NewImpala() {
	mu.Lock()
	defer mu.Unlock()
	var err error
	impala, err = impalathing.Connect(config.Configs.Impala.Host, config.Configs.Impala.Port)
	if err != nil {
		panic(fmt.Sprintf("Impala connect err:%s", err))
	}
}

func CloseImpala() {
	_ = impala.Close()
}

func query(q string) (rs impalathing.RowSet, err error) {
	rs, err = impala.Query(q)
	if err != nil {
		return rs, err
	}
	status, err := rs.Wait()
	if err != nil {
		return rs, err
	}
	if !status.IsSuccess() {
		return rs, fmt.Errorf("unsuccessful query execution: %+v", status)
	}
	return
}

func QueryAll(q string) ([]map[string]interface{}, error) {
	mu.Lock()
	defer mu.Unlock()
	res, err := query(q)
	if err != nil {
		return nil, err
	}
	return res.FetchAll(), nil
}

func Query(q string) error {
	_, err := QueryAll(q)
	return err
}

func FetchWithoutKey(q string) ([]interface{}, error) {
	res, err := QueryAll(q)
	if err != nil {
		return nil, err
	}
	response := make([]interface{}, 0)
	for _, i2 := range res {
		row := make([]interface{}, len(i2))
		i := 0
		for _, i3 := range i2 {
			row[i] = i3
			i++
		}
		response = append(response, row)
	}

	return response, nil
}
