package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

// 解析ini文件

// mysql配置对应的结构体
type MySQL struct {
	UserName string `mysql:"username"`
	Password string `mysql:"password"`
	Addr string `mysql:"addr"`
	Port string `mysql:"port"`
}

// redis配置对应的结构体
type Redis struct {
	Host string `redis:"host"`
	Port string `redis:"port"`
}

// 解析ini文件
func parseIni(fileName string, mysql *MySQL, redis *Redis) (err error) {
	// 传入的参数不是指针则返回
	mysqlT := reflect.TypeOf(mysql)
	mysqlV := reflect.ValueOf(mysql)
	redisT := reflect.TypeOf(redis)
	redisV := reflect.ValueOf(redis)
	if mysqlT.Kind() != reflect.Ptr || redisT.Kind() != reflect.Ptr {
		return errors.New("传入的不是指针参数")
	}

	// 传入的参数指针不是指向结构体则返回
	if mysqlT.Elem().Kind() != reflect.Struct || redisT.Elem().Kind() != reflect.Struct {
		return errors.New("参数应指向结构体")
	}

	// 打开文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.New("open file err")
	}

	// 将文件内容全部转换为字符串
	lines := strings.Split(string(data), "\r\n")

	// 逐行判断是否满足格式
	for idx := 0; idx < len(lines); {
		line := lines[idx]
		line = strings.Trim(line, " \r\n")

		// 空行或注释则直接跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || len(line) == 0 {
			idx++
			continue
		}

		// 甄别格式错误的line，直接返回错误(key = val)
		// 1. 标题的格式 [title]
		isTile, title, err := isTitle(line, idx)
		if err != nil {
			return err
		}
		if isTile { // 找到标题
			// 从标题的下一行开始到下一个title结束作为该title的body
			idx++ // 跳到下一行
			for { // 循环判断是否是标题，如果不是标题则判断是否是bodyLine
				if idx >= len(lines) {
					break
				}
				isTile, _, err = isTitle(lines[idx], idx)
				if !isTile {
					// 空行或注释则直接跳过
					if strings.HasPrefix(lines[idx], ";") || strings.HasPrefix(lines[idx], "#") || len(lines[idx]) == 0 {
						idx++
						continue
					}
					isBodyLine, err := isBody(lines[idx], idx) // 是否是bodyLine
					if err != nil {
						return err
					}
					if isBodyLine {
						// 解析bodyLine
						bodyLine := strings.Trim(lines[idx], " \r\n")
						tag := strings.ToLower(strings.TrimSpace(bodyLine[: strings.Index(bodyLine, "=")]))
						val := strings.TrimSpace(bodyLine[strings.Index(bodyLine, "=") + 1 :])
						switch title {
						case "mysql":
							for i := 0; i < mysqlT.Elem().NumField(); i++ {
								if mysqlT.Elem().Field(i).Tag.Get("mysql") == tag {
									mysqlV.Elem().Field(i).SetString(val)
								}
							}
						case "redis":
							for i := 0; i < redisT.Elem().NumField(); i++ {
								if redisT.Elem().Field(i).Tag.Get("redis") == tag {
									redisV.Elem().Field(i).SetString(val)
								}
							}
						}
					}
					idx++
				} else {
					break // 如果是title 不能跳过，要继续从这一行做判断
				}
			}
		}
	}
	return
}

func isBody(line string, idx int) (res bool, err error) {
	if !strings.Contains(line,"=") {
		return false, fmt.Errorf("格式错误，at line: %d", idx)
	} else {
		tag := strings.TrimSpace(line[:strings.Index(line, "=")])
		if len(tag) == 0 {
			return false, fmt.Errorf("格式错误，at line: %d", idx)
		}
		return true, nil
	}
}

func isTitle(line string, idx int) (res bool, title string, err error) {
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		line = line[strings.Index(line,"[") + 1 : strings.Index(line,"]")]
		line = strings.TrimSpace(line)
		if strings.ToLower(line) != "mysql" && strings.ToLower(line) != "redis" {
			return false, "", fmt.Errorf("标题格式错误，at line: %d", idx)
		}
		return true, line, nil
	}
	return false, "", nil
}

func main() {
	fileName := "E:/GoWorkPlace/src/myLog/ini/xxx.ini"
	mySQLConfig := &MySQL{}
	redisConfig := &Redis{}
	err := parseIni(fileName, mySQLConfig, redisConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s %s\n", mySQLConfig.Addr, redisConfig.Host)
}

