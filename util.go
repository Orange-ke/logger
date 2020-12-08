package myLog

import (
	"path"
	"runtime"
)

func GetCallerInfo(skip int) (fileName, funcName string, line int) {
	pc, fileName, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}
	// 根据pc拿到当前执行的函数名
	funcName = runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)
	// 从文件全路径剥离出文件名
	fileName = path.Base(fileName)
	return
}