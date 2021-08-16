package log

import (
	"fmt"
	"log"
	"os"
	//"github.com/tddhit/tools/goid"
)

const (
	TRACE = 1 + iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
	PANIC
)

var flag = log.LstdFlags | log.Lshortfile | log.Lmicroseconds
var logger = log.New(os.Stderr, "", flag)
var logLevel = TRACE
var logPath string

func Init(path string, level int) {
	if path != "" {
		logPath = path
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			Error("failed open file: %s, %s", path, err)
		} else {
			logger = log.New(file, "", flag)
		}
	}
	if level >= DEBUG && level <= PANIC {
		logLevel = level
	}
}

func Reopen() {
	Init(logPath, logLevel)
}

func Panicf(format string, v ...interface{}) {
	if logLevel <= PANIC {
		format = "[PANIC] " + format
		s := format + fmt.Sprintf(format, v...)
		logger.Output(2, s)
		panic(s)
	}
}

func Panic(v ...interface{}) {
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
		format = fmt.Sprintf("[INFO] GID(%d) ", 0) + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if logLevel <= INFO {
		s := fmt.Sprintf("[INFO] GID(%d) ", 0) + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		format = fmt.Sprintf("[DEBUG] GID(%d) ", 0) + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		s := fmt.Sprintf("[DEBUG] GID(%d) ", 0) + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func Tracef(calldepth int, format string, v ...interface{}) {
	if logLevel <= TRACE {
		format = fmt.Sprintf("[Trace] GID(%d) ", 0) + format
		logger.Output(calldepth, fmt.Sprintf(format, v...))
	}
}

func Trace(calldepth int, v ...interface{}) {
	if logLevel <= TRACE {
		s := fmt.Sprintf("[Trace] GID(%d) ", 0) + fmt.Sprintln(v...)
		logger.Output(calldepth, s)
	}
}
