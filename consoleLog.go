package myLog

// 往终端中写日志

import (
	"fmt"
	"os"
	"time"
)

// 往终端中记录日志的结构体
type consoleLogger struct {
	level Level // 大于这个级别的日志才记录（区分生产环境和开发环境）
	// ...
}

// 获取consoleLogger对象的工厂方法构造函数
func NewConsoleLogger(levelStr string) *consoleLogger {
	level := ParseLogLevel(levelStr)
	clObj := &consoleLogger{
		level: level,
	}
	return clObj
}

// 方法
// Debug debug 方法
func (f *consoleLogger) Debug(format string, args ...interface{}) {
	f.log(DEBUG, format, args...)
}

// Info info方法
func (f *consoleLogger) Info(format string, args ...interface{}) {
	f.log(INFO, format, args...)
}

// Warn warn方法
func (f *consoleLogger) Warn(format string, args ...interface{}) {
	f.log(WARN, format, args...)
}

// Error error方法
func (f *consoleLogger) Error(format string, args ...interface{}) {
	f.log(ERROR, format, args...)
}

// Fatal fatal方法
func (f *consoleLogger) Fatal(format string, args ...interface{}) {
	f.log(FATAL, format, args...)
}

// Close 终端标准输出不需要关闭
func (f *consoleLogger) Close() {

}

// 记录日志
// 将公用的记录日志的功能封装成一个单独的方法
func (f *consoleLogger) log(level Level, format string, args ...interface{}) {
	if f.level > level { // 一般上线后不会打印DEBUG信息
		return
	}
	nowStr := time.Now().Format("2006-01-02 15:04:05:05.000")
	funcName, fileName, line := GetCallerInfo(3)
	// 拼接日志格式
	// 格式：[时间][文件：行号][函数名][日志级别] 日志信息
	logStr := fmt.Sprintf("[%s] [%s : %d] [%s] [%s] %s \n",
		nowStr, fileName, line, funcName, GetLevelStr(f.level), format)

	_, err := fmt.Fprintf(os.Stdout, logStr, args...)
	if err != nil {
		panic(fmt.Sprintf("Write to console err: %v \n", err))
	}
}
