package main

import (
	"myLog"
	"time"
)

var logger myLog.Logger

func main() {
	logger = myLog.NewFileLogger("Debug", "./", "test.log")
	defer logger.Close()
	// f1 := myLog.NewConsoleLogger("Info")
	userId := 100
	start := time.Now().Unix()
	for {
		end := time.Now().Unix()
		logger.Debug("用户id %d 一直尝试登录", userId)
		logger.Info("Info 测试")
		logger.Error("Error 测试")
		if end - start > 10 {
			break
		}
	}
	logger.CloseChan()
	for {
		_, ok := logger.Exit()
		if !ok {
			break
		}
	}
}
