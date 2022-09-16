package log

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const baseFormat = "2006-01-02 15:04:05"

type Config struct {
	Level    string `json:",default=info"`
	FileName string `json:",optional"`
	MaxSize  int    `json:",default=500"`
	MaxAge   int    `json:",default=3"`
}

func init() {
	c := Config{}
	SetUp(&c)
}

func SetUp(c *Config) {
	// set log json formatter
	logrus.SetFormatter(&SacreFormatter{})

	// set log level
	level := logrus.InfoLevel
	switch strings.ToLower(c.Level) {
	case "debug":
		level = logrus.DebugLevel
	case "error":
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)

	// set log policy for file
	if c.FileName != "" && strings.HasPrefix(c.FileName, "/") {
		logrus.SetOutput(&lumberjack.Logger{
			Filename: c.FileName,
			MaxSize:  c.MaxSize,
			MaxAge:   c.MaxAge,
		})
	}
}

type SacreFormatter struct{}

func (my *SacreFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	t := entry.Time.Format(baseFormat)
	message := entry.Message[1 : len(entry.Message)-1]
	contents := []byte(fmt.Sprintf("[%s] [%s] [%s] [RoutineID=%s] %s\n", level, t, getFileLine(), getGoroutinedID(), message))
	return contents, nil
}

func Debug(args ...interface{}) {
	logrus.Debug(args)
}

func Info(args ...interface{}) {
	logrus.Info(args)
}

func Warn(args ...interface{}) {
	logrus.Warn(args)
}

func Error(args ...interface{}) {
	logrus.Error(args)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args)
}

func ServerError(args ...interface{}) {
	logrus.Error(args, string(debug.Stack()))
}

func ParseLevel(level string) (logrus.Level, error) {
	return logrus.ParseLevel(level)
}

func getFileLine() string {
	fl := ""
	if _, file, line, ok := runtime.Caller(9); ok {
		fl = fmt.Sprintf("%s:%d", file, line)
	}
	return fl
}

func getGoroutinedID() string {
	buf := [64]byte{}
	n := runtime.Stack(buf[:], false)
	stk := strings.TrimPrefix(string(buf[:n]), "goroutine")
	firstField := strings.Fields(stk)[0]
	return firstField
}
