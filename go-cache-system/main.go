/**
 * @Author: djh
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/9/23 14:41
 */

package main

import (
	"cache-system/config"
	"cache-system/scan"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func beforeAction() {
	config.NewConfig()
}

func main() {
	beforeAction()
	var finishUp = make(chan struct{})
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, os.Interrupt, os.Kill)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("Caught sig: %+v", sig)
		finishUp <- struct{}{}
	}()

guard:
	for {
		select {
		case <-finishUp:
			fmt.Println("Stopping progress")
			break guard
		default:
			err := scan.Run()
			if err != nil {
				break guard
			}
			scan.Log.Info("sleep....")
			time.Sleep(time.Duration(config.Configs.Select.Sleep) * time.Second)

		}
	}
}
