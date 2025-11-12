package logger

import (
	"fmt"
	"log"
	"strings"
)

type logLevel string

const (
	loggerLevelDebug logLevel = "DBG"
	loggerLevelWarn  logLevel = "WRN"
	loggerLevelError logLevel = "ERR"
	loggerLevelFatal logLevel = "FTL"
)

func Debug(msg string, keyvals ...any) {
	logPrint(loggerLevelDebug, msg, keyvals...)
}

func Warn(msg string, keyvals ...any) {
	logPrint(loggerLevelWarn, msg, keyvals...)
}

func Error(msg string, keyvals ...any) {
	logPrint(loggerLevelError, msg, keyvals...)
}

func Fatal(msg string, keyvals ...any) {
	logPrint(loggerLevelFatal, msg, keyvals...)
}

// ParseArgs  解析error和message用统一格式展示出来
func ParseArgs(kvs ...any) (string, map[string]any, error) {
	countArgs := len(kvs)
	// 如果可变参数个数为0，肯定没有error
	if countArgs == 0 {
		return "", nil, nil
	}

	var err error
	var msgValue string
	args := make(map[string]any)
	for i := 0; i < countArgs-1; i = i + 2 {
		// 第i个值作为map的key, 第i+1个值作为map的value
		k, ok := kvs[i].(string)
		if !ok {
			continue
		}

		switch strings.ToLower(k) {
		case "level", "caller":
			// do nothing
		case "err", "panic": // err and panic must be panic type
			switch v := kvs[i+1].(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", kvs[i+1])
			}
		case "msg", "message":
			msgValue = fmt.Sprintf("%v", kvs[i+1])
		default:
			args[k] = kvs[i+1]
		}
	}
	return msgValue, args, err
}

// logPrint log structure message and key values
func logPrint(level logLevel, msg string, kvs ...any) {
	_, fields, err := ParseArgs(kvs...)

	outputs := make([]string, 0)
	for k, v := range fields {
		outputs = append(outputs, fmt.Sprintf("%s=\"%v\"", k, v))
	}

	logFn := log.Printf
	if level == loggerLevelFatal {
		logFn = log.Fatalf
	}

	if len(outputs) > 0 {
		if err != nil {
			logFn("%s msg=\"%s\" %s error=\"%v\"", level, msg, strings.Join(outputs, " "), err)
		} else {
			logFn("%s msg=\"%s\" %s", level, msg, strings.Join(outputs, " "))
		}
	} else {
		if err != nil {
			logFn("%s msg=\"%s\" error=\"%v\"", level, msg, err)
		} else {
			logFn("%s msg=\"%s\"", level, msg)
		}
	}
}
