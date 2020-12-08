package main

// 利用反射将配置文件中的参数赋值给对应的结构体

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// config是一个日志配置项
type Config struct {
	FilePath string `conf:"file_path"`
	FileName string `conf:"file_name"`
	MaxSize int64 `conf:"max_size"`
}

// 从conf文件中读取内容并赋值给结构体指针
func parseConf(fileName string, config interface{}) (err error) {

	// config 参数必须是prt类型，否则无法赋值
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		err = errors.New("config的类型必须是指针类型")
		return err
	}

	// *config 必须是结构体
	if reflect.TypeOf(config).Elem().Kind() != reflect.Struct {
		err = errors.New("config必须是结构体的指针")
		return err
	}

	// 获取config的类型和值
	confType := reflect.TypeOf(config)
	configVal := reflect.ValueOf(config)

	// 一次性读取配置文件的所有内容
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("读取文件失败\n")
		return err
	}

	// 按照换行符分割为 字符串 切片
	strSlice := strings.Split(string(data), "\r\n")

	// 逐行解析
	for idx, line := range strSlice {
		line := strings.TrimSpace(line)
		// 如果为空行或者以#开头，则跳过
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// 如果遇到不合格的描述也要跳过 即 不含=
		if strings.Index(line, "=") == -1 {
			err = fmt.Errorf("配置文件%d行有错误", idx + 1)
			return err
		}

		// 将=左右的内容截取，作为tag 和 val的取值
		tag := strings.TrimSpace(line[:strings.Index(line, "=")])
		val := strings.TrimSpace(line[strings.Index(line, "=") + 1:])

		// 如果 tag为空也非法
		if len(tag) == 0 {
			err = fmt.Errorf("配置文件%d行有错误", idx + 1)
			return err
		}

		// 循环结构体的字段查找匹配的tag，然后根据字段的类型选择合适的方法设定值
		for i := 0; i < confType.Elem().NumField(); i++ {
			if confType.Elem().Field(i).Tag.Get("conf") == tag {
				switch confType.Elem().Field(i).Type.Kind() {
				case reflect.String:
					configVal.Elem().Field(i).SetString(val)
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
					intVal, err := strconv.ParseInt(val,10,64)
					if err != nil {
						return err
					}
					configVal.Elem().Field(i).SetInt(intVal)
				default:
					break
				}
				break
			}
		}
	}
	return
}

func main() {
	// 1. 打开文件
	// 2. 读取内容
	// 3. 一行一行的读取内容，根据tag找结构体里面对应的字段
	// 4. 找到则赋值

	var c = &Config{}
	fmt.Println(c)
	err := parseConf("xxx.conf", c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", *c)
}
