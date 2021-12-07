package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"runtime"
)

type Config struct {
	Output     string
	MaxSize    int //日志文件最大值 单位兆,默认100兆
	MaxAge     int //日志保留最大天数，默认不删除
	MaxBackups int //要保留的旧日志文件的最大数量,默认都保留
	Level      string
}

func Init(c *Config) error {
	lvl, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return err
	}
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   c.Output,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge, //days
	})
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	return nil
}

func ErrLog(info string, err error) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	var name, file string
	var line int
	if f != nil {
		name = f.Name()
		file, line = f.FileLine(pc[0])
	}
	logrus.WithFields(logrus.Fields{
		"name": name,
		"path": fmt.Sprintf("%s %d", file, line),
		"err":  err,
	}).Error(info)
}

func Infof(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Printf("\n")
	logrus.Infof(str, args...)
}

func Debugf(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Printf("\n")
	logrus.Debugf(str, args...)
}

func Fatalf(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	fmt.Printf("\n")
	logrus.Fatalf(str, args...)
}
