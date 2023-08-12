package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

type LoggingLevel int

const (
	DEBUG   LoggingLevel = 0
	GENERAL LoggingLevel = 1
	PROD    LoggingLevel = 2
)

var loggingLevel LoggingLevel = 0

func SetUpLogging(l LoggingLevel) {
	loggingLevel = l
}

const timeLayout = "[2006/01/02 @ 15:04:05] "

var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgHiYellow).SprintFunc()
var blue = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func currentTimeString() string {
	return time.Now().Format(timeLayout)
}

func SUCCESS(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= PROD {
		fmt.Println(green(currentTimeString()+"[SUCCESS/"+scope+"] ") + s)
	}
}

func INFO(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= GENERAL {
		fmt.Println(blue(currentTimeString()+"[INFO/"+scope+"] ") + s)
	}
}

func LOGGING(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= DEBUG {
		fmt.Println(currentTimeString() + "[LOGGING/" + scope + "] " + s)
	}
}

func WARNING(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= GENERAL {
		fmt.Println(yellow(currentTimeString()+"[WARNING/"+scope+"] ") + s)
	}
}

func ERROR(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= PROD {
		fmt.Println(red(currentTimeString()+"[ERROR/"+scope+"]\n") + s)
	}
}

func FATAL(s string, scope string) {
	if scope == "" {
		scope = "global"
	}
	if loggingLevel <= PROD {
		log.Fatal(red(currentTimeString()+"[FATAL/"+scope+"]\n") + s)
	}
}
