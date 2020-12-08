package myLog

import "strings"

// 自定义一个日志库，实现日志记录的功能

// 日志分级别
// DEBUG TRACE INFO WARN ERROR FATAL
type Level uint16

// 具体的日志级别常量
const (
	DEBUG Level = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

// 定义一个logger接口
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Close()
	CloseChan()
	Exit() (bool, bool)
}

// 日志记录结构体
type LogData struct {
	msg string
	logLevel Level
	logTime string
	funcName string
	fileName string
	lineNum int
	args []interface{}
}

// 根据对应的数字获取对应的字符串描述
func GetLevelStr (level Level) string {
	switch level {
	case 0:
		return "DEBUG"
	case 1:
		return "TRACE"
	case 2:
		return "INFO"
	case 3:
		return "WARN"
	case 4:
		return "ERROR"
	case 5:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

// 根据用户输入的字符串转化为对应的Level
func ParseLogLevel(levelStr string) Level {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return DEBUG
	case "trace":
		return TRACE
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return DEBUG
	}
}
