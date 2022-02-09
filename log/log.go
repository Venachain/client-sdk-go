package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gggwvg/logrotate"
	"github.com/sirupsen/logrus"
)

type Level uint32
type LogFormatter struct{}

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

var logrusPackage = "github.com/sirupsen/logrus"
var logPackage = "github.com/PlatONEnetwork/PlatONE-Cross/log"

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type LogConfig struct {
	Level Level  `yaml:"level"`
	Size  string `yaml:"size"`
	Path  string `yaml:"path"`
}

var DefaultLogConfi = LogConfig{
	5,
	"30m",
	"./log/data/",
}

func init() {
	err := MakeLogConfig(DefaultLogConfi)
	if err != nil {
		return
	}
}

func MakeLogConfig(cfg LogConfig) error {
	logPath := cfg.Path
	if logPath == "" {
		return fmt.Errorf("unable to parse the log path from config")
	}
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir failed![%v]\n", err)
	}
	time := time.Now()
	fileDate := time.Format("2006-01-02")
	filename := fmt.Sprintf("%s%s.log", logPath, fileDate)
	opts := []logrotate.Option{
		logrotate.File(filename),
	}
	if cfg.Size == "" {
		return fmt.Errorf("unable to parse the log size from config")
	}
	opts = append(opts, logrotate.RotateSize(cfg.Size))
	logger, err := logrotate.NewLogger(opts...)
	if err != nil {
		return fmt.Errorf("Newlogger is error[%v]\n", err)
	}
	level := cfg.Level
	logrus.SetLevel(logrus.Level(level))
	logrus.SetOutput(os.Stdout)
	logrus.SetOutput(logger)
	writers := []io.Writer{
		logger,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to log to file")
	}
	logrus.SetFormatter(new(LogFormatter))
	logger.Close()
	return nil
}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var len int

	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	var caller *runtime.Frame
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != logPackage {
			caller = &f //nolint:scopelint
			break
		}
	}
	if caller != nil {
		file = filepath.Base(caller.File)
		len = caller.Line
	}
	msg := fmt.Sprintf("%s [%s:%d][%s] %s\n", timestamp, file, len, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func Info(format string, a ...interface{}) {
	logrus.Infof(format, a...)
}

func Warn(format string, a ...interface{}) {
	logrus.Warnf(format, a...)
}

func Debug(format string, a ...interface{}) {
	logrus.Debugf(format, a...)
}

func Error(format string, a ...interface{}) {
	logrus.Errorf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	logrus.Fatalf(format, a...)
}

func Panic(format string, a ...interface{}) {
	logrus.Panicf(format, a...)
}
