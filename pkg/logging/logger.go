package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	LogLevels = map[string]logrus.Level{
		"warn":  3,
		"info":  4,
		"debug": 5,
		"trace": 6,
	}
)

type Hooks struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *Hooks) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hooks) Levels() []logrus.Level {
	return h.LogLevels
}

var (
	entry *logrus.Entry
	once  sync.Once
)

type Logger struct {
	*logrus.Entry
}

func GetLogger(logLevel string) *Logger {
	initialize(logLevel)
	return &Logger{Entry: entry}
}

func initialize(logLevel string) {
	once.Do(func() {
		logger := logrus.New()

		logger.SetReportCaller(true)
		logger.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				runes := []rune(frame.File)

				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}

				var thirdSlashId int

				for i, count := 0, 0; i < len(runes)-1; i++ {
					if count == 3 {
						thirdSlashId = i
						break
					}
					if runes[i] == '/' {
						count++
					}
				}

				filename := frame.File[len(runes)-thirdSlashId:]

				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
			DisableColors:   false,
			TimestampFormat: time.RFC1123,
		}

		err := os.Mkdir("logs", os.ModePerm)
		if err != nil && !os.IsExist(err) {
			logrus.Fatalf("error while creating dir logs: %s", err.Error())
		}

		logFile, err := os.OpenFile("logs/log_file.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			logrus.Fatalf("error while creating custom logging: %s", err.Error())
		}

		logger.SetOutput(io.Discard)

		logger.AddHook(&Hooks{
			Writer:    []io.Writer{logFile, os.Stdout},
			LogLevels: logrus.AllLevels,
		})

		logger.SetLevel(LogLevels[logLevel])
		fmt.Println(logger.Level, LogLevels[logLevel])

		entry = logrus.NewEntry(logger)

		fmt.Println(entry.Logger.Level)
	})
}
