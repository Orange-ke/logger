package myLog

// 往文件中写日志

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件中记录日志的结构体
type fileLogger struct {
	level Level // 大于这个级别的日志才记录（区分生产环境和开发环境）
	logFilePath string
	logFileName string
	logFile *os.File
	errorFile *os.File
	max int64
	splitTime *time.Time // 上一次切分的时间
	logChan chan *LogData
	exitChan chan bool
}

// 获取fileLogger对象的工厂方法构造函数
func NewFileLogger(levelStr string, logFilePath, logFileName string) *fileLogger{
	level := ParseLogLevel(levelStr)
	flObj := &fileLogger{
		level: level,
		logFilePath: logFilePath,
		logFileName: logFileName,
		max: 10 * 1024 * 1024,
		logChan: make(chan *LogData, 50000), // 初始化一个通道把日志文件写入
		exitChan: make(chan bool, 1),
	}
	flObj.initFileLogger()
	return flObj
}

// 初始化文件日志的文件句柄
func (f *fileLogger) initFileLogger() {
	filePath := path.Join(f.logFilePath, f.logFileName)
	// 打开文件
	file, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("open file: %s failed, %v \n", filePath, err))
	}
	f.logFile = file
	// 打开记录错误日志文件
 	errFile, err := os.OpenFile(fmt.Sprintf("%s.err", filePath), os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("open file: %s failed, %v \n", filePath, err))
	}
	f.errorFile = errFile
	// 初始化后启动协程从管道中读取数据
	go f.writeLogBackground()
}

// 方法
// Debug debug 方法
func (f *fileLogger) Debug(format string, args ...interface{}) {
	f.log(DEBUG, format, args...)
}

// Info info方法
func (f *fileLogger) Info(format string, args ...interface{}) {
	f.log(INFO, format, args...)
}

// Warn warn方法
func (f *fileLogger) Warn(format string, args ...interface{}) {
	f.log(WARN, format, args...)
}

// Error error方法
func (f *fileLogger) Error(format string, args ...interface{}) {
	f.log(ERROR, format, args...)
}

// Fatal fatal方法
func (f *fileLogger) Fatal(format string, args ...interface{}) {
	f.log(FATAL, format, args...)
}

// Close close关闭日志文件
func (f *fileLogger) Close() {
	_ = f.logFile.Close()
	_ = f.errorFile.Close()
}

func (f *fileLogger) CloseChan() {
	close(f.logChan)
}

func (f *fileLogger) Exit() (bool, bool) {
	v, ok := <-f.exitChan
	return v, ok
}

// 记录日志
// 将公用的记录日志的功能封装成一个单独的方法
func (f *fileLogger) log(level Level, format string, args ...interface{}) {
	if f.level > level { // 一般上线后不会打印DEBUG信息
		return
	}
	nowStr := time.Now().Format("2006-01-02 15:04:05:05.000")
	funcName, fileName, line := GetCallerInfo(3)
	// 构造logData结构体
	logData := &LogData{
		msg: format,
		logLevel: level,
		logTime: nowStr,
		funcName: funcName,
		fileName: fileName,
		lineNum: line,
		args: args,
	}
	// 放入通道
	select {
	case f.logChan <- logData:
	default: // 丢弃日志信息
		// 取出之前的日志，放入新的日志
		//<- f.logChan
		//f.logChan <- logData
	}
}

func (f *fileLogger) writeLogBackground() {
	for {
		logData, ok := <- f.logChan
		if !ok {
			break
		}
		// 拼接日志格式
		// 格式：[时间][文件：行号][函数名][日志级别] 日志信息
		logStr := fmt.Sprintf("[%s] [%s : %d] [%s] [%s] %s \n",
			logData.logTime, logData.funcName, logData.lineNum, logData.fileName, GetLevelStr(logData.logLevel), logData.msg)
		// 判断是否要切分文件
		// 按文件大小切分
		if f.checkFileSize(f.logFile) {
			f.logFile = f.splitLogFile(f.logFile)
		}
		_, err := fmt.Fprintf(f.logFile, logStr, logData.args...)
		if err != nil {
			panic(fmt.Sprintf("Log write to %s err", f.logFilePath + f.logFileName))
		}
		// 如果是error或这fatal级别的日志还要记录到发f.errFile中
		if logData.logLevel >= ERROR {
			// 判断是否要切分文件
			// 按照大小切分
			if f.checkFileSize(f.logFile) {
				f.errorFile = f.splitLogFile(f.errorFile)
			}
			_, err = fmt.Fprintf(f.errorFile, logStr, logData.args...)
			if err != nil {
				panic(fmt.Sprintf("Log write to %s err \n", f.logFilePath + ".err" + f.logFileName))
			}
		}
	}
	close(f.exitChan)
}

// 判断日志文件是否超过了maxSize
func (f *fileLogger) checkFileSize(file *os.File) bool {
	// 往文件里写日志之前，要做一个检查，判断当前日志文件的大小是否超过maxSize
	fileInfo , _ := file.Stat()
	fileSize := fileInfo.Size()
	return fileSize >= f.max
}

// 按大小切分文件
func (f *fileLogger) splitLogFile(file *os.File) *os.File {
	// 切分文件
	// 1. 把原来的文件关闭
	_ = file.Close()
	// 2. 备份原来的文件
	fileName := file.Name()
	backFileName := fmt.Sprintf("%s_%v.back", fileName, time.Now().Unix())
	err := os.Rename(fileName, backFileName)
	if err != nil {
		panic(fmt.Errorf("rename err: %v", err))
	}
	// 3. 新建一个文件
	fileObj, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("open file err: %v", err))
	}
	return fileObj
}

// 判断是否需要按照时间切分
func (f *fileLogger) checkSplitByTime(sec int64) bool {
	start := time.Now()
	if f.splitTime == nil {
		f.splitTime = &start
		return false
	}
	duration := start.Unix() - f.splitTime.Unix()
	if duration >= sec {
		f.splitTime = &start
		return true
	}
	return false
}
