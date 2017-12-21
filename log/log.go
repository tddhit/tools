package log

import (
	"fmt"
	"log"
	"os"
)

const (
	DEBUG = 1 + iota
	INFO
	WARNING
	ERROR
	FATAL
	PANIC
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
var logLevel = INFO

func InitLogger(logPath string) {
	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			Error("failed open file: %s, %s", logPath, err)
		} else {
			logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		}
	}
	if logLevel >= DEBUG && logLevel <= PANIC {
		logLevel = logLevel
	}
}

func Panicf(format string, v ...interface{}) {
	if logLevel <= PANIC {
		format = "[PANIC] " + format
		s := format + fmt.Sprintf(format, v...)
		logger.Output(2, s)
		panic(s)
	}
}

func Panic(format string, v ...interface{}) {
	if logLevel <= PANIC {
		s := "[PANIC] " + fmt.Sprintln(v...)
		logger.Output(2, s)
		panic(s)
	}
}

func Fatalf(format string, v ...interface{}) {
	if logLevel <= FATAL {
		format = "[FATAL] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

func Fatal(v ...interface{}) {
	if logLevel <= FATAL {
		s := "[FATAL] " + fmt.Sprintln(v...)
		logger.Output(2, s)
		os.Exit(1)
	}
}

func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		format = "[ERROR] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	if logLevel <= ERROR {
		s := "[ERROR] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func Warnf(format string, v ...interface{}) {
	if logLevel <= WARNING {
		format = "[WARNING] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(v ...interface{}) {
	if logLevel <= WARNING {
		s := "[WARNING] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		format = "[INFO] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if logLevel <= INFO {
		s := "[INFO] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		format = "[DEBUG] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		s := "[DEBUG] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}
