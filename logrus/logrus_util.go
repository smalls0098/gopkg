package logrus

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Prod        bool
	Filename    string
	MaxSize     int32
	MaxBackups  int32
	MaxAge      int32
	Compress    bool
	ProdShowStd bool
}

func NewLogrus(config Config) *logrus.Logger {
	logs := logrus.New()
	if config.Prod {
		logs.SetLevel(logrus.InfoLevel)
		logs.SetFormatter(&logrus.TextFormatter{})
		ljLogger := &lumberjack.Logger{
			Filename:   config.Filename,        //指定日志存储位置
			MaxSize:    int(config.MaxSize),    //日志的最大大小（M）
			MaxBackups: int(config.MaxBackups), //日志的最大保存数量
			MaxAge:     int(config.MaxAge),     //日志文件存储最大天数
			Compress:   config.Compress,        //是否执行压缩
			LocalTime:  true,                   //使用本地时间命名 默认为utc时间命名
		}
		if config.ProdShowStd {
			logs.SetOutput(io.MultiWriter(ljLogger, os.Stderr))
		} else {
			logs.SetOutput(io.MultiWriter(ljLogger))
		}
	} else {
		logs.SetLevel(logrus.TraceLevel)
		logs.SetOutput(os.Stderr)
	}
	return logs
}
