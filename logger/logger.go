package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/lestrrat/go-file-rotatelogs"
	"path"
	"strings"
	"io"
	"os"
)

type Logger struct {
	Level string `json:"level" description:"日志等级"`
	WriterMap struct {
		Debug string `json:"debug"`
		Info  string `json:"info"`
		Warn  string `json:"warn"`
		Error string `json:"error"`
	} `json:"writerMap"`
	File struct {
		Path             string `json:"path" description:"日志文件位置"`
		MaxAgeHour       int    `json:"maxAgeHour" description:"日志文件保存小时时间"`
		RotationTimeHour int    `json:"rotationTimeHour" description:"日志文件分割时间"`
	} `json:"file"`
}

var DefaultConfig = `
{
    "level": "debug",
    "writerMap": {
      "debug": "Console",
      "info": "Console",
      "warn": "Console",
      "error": "Console"
    },
    "file": {
      "path": "D:/Java/go_repository/src/per.zdh.org/cmdb/logs",
      "maxAgeHour": 50,
      "rotationTimeHour": 24
    }
}
`

// 项目初始化引入配置信息
func New(option Logger) error {
	writerMap := lfshook.WriterMap{}
	var fileWriter io.Writer

	if &option.File != nil {

		_, statErr := os.Stat(option.File.Path)
		// 如果目录不存在则创建该目录
		if os.IsNotExist(statErr) {
			if err := os.MkdirAll(option.File.Path, 0755); err != nil {
				return err
			}
		}

		writer, err := rotatelogs.New(
			path.Join(option.File.Path, "%Y%m%d%H%M.log"),
			rotatelogs.WithMaxAge(time.Hour*time.Duration(option.File.MaxAgeHour)),             // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Hour*time.Duration(option.File.RotationTimeHour)), // 日志切割时间间隔
		)
		if err != nil {
			return err
		}
		fileWriter = writer
	}

	switch strings.ToUpper(option.WriterMap.Debug) {
	case "CONSOLE":
		writerMap[logrus.DebugLevel] = os.Stdout
	case "FILE":
		writerMap[logrus.DebugLevel] = fileWriter
	default:
		writerMap[logrus.DebugLevel] = os.Stdout
	}

	switch strings.ToUpper(option.WriterMap.Info) {
	case "CONSOLE":
		writerMap[logrus.InfoLevel] = os.Stdout
	case "FILE":
		writerMap[logrus.InfoLevel] = fileWriter
	default:
		writerMap[logrus.InfoLevel] = os.Stdout
	}

	switch strings.ToUpper(option.WriterMap.Warn) {
	case "CONSOLE":
		writerMap[logrus.WarnLevel] = os.Stdout
	case "FILE":
		writerMap[logrus.WarnLevel] = fileWriter
	default:
		writerMap[logrus.WarnLevel] = os.Stdout
	}

	switch strings.ToUpper(option.WriterMap.Error) {
	case "CONSOLE":
		writerMap[logrus.ErrorLevel] = os.Stdout
	case "FILE":
		writerMap[logrus.ErrorLevel] = fileWriter

	default:
		writerMap[logrus.ErrorLevel] = os.Stdout
	}

	// 为不同级别设置不同的输出目的
	lfHook := lfshook.NewHook(writerMap, &logrus.JSONFormatter{})
	switch strings.ToUpper(option.Level) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
	logrus.AddHook(lfHook)
	return nil
}
