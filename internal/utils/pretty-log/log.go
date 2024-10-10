package prettylog

import (
	"fmt"
	"log"
	"strconv"
)

const (
	timeFormat = "[15:04:05.000]"

	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var logLevel = INFO

func SetLevel(level int) {
	logLevel = level
}

func Debug(msg string, data ...any) {
	if logLevel > DEBUG {
		return
	}
	log.Println(colorize(green, msg, data...))
}

func Info(msg string, data ...any) {
	if logLevel > INFO {
		return
	}
	log.Println(colorize(blue, msg, data...))
}

func Warning(msg string, data ...any) {
	if logLevel > WARN {
		return
	}
	log.Println(colorize(yellow, msg, data...))
}

func Error(msg string, data ...any) {
	log.Println(colorize(red, msg, data...))
}

func colorize(colorCode int, msg string, data ...any) string {
	fmtMsg := fmt.Sprintf(msg, data...)
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), fmtMsg, reset)
}
