package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	infoLog = log.New(os.Stdout,
		"PUFCTL - INFO: ",
		log.Ldate|log.Ltime,
	)
	warnLog = log.New(os.Stderr,
		"PUFCTL - Warn: ",
		log.Ldate|log.Ltime,
	)
	errorLog = log.New(os.Stderr,
		"PUFCTL - ERROR: ",
		log.Ldate|log.Ltime,
	)
	debugLog = log.New(os.Stdout,
		"PUFCTL - DEBUG: ",
		log.Ldate|log.Ltime,
	)
)

// Infoln logs to stdout in the format of fmt.Println
func Infoln(msg ...interface{}) {
	infoLog.Println(msg...)
}

// Warnln logs to stderr in the format of fmt.Println
func Warnln(msg ...interface{}) {
	warnLog.Println(msg...)
}

// Errorln logs a string to stderr and exits the program with exit code 1 in the format of fmt.Fatalln
func Errorln(msg ...interface{}) {
	errorLog.Fatalln(msg...)
}

// Debugln logs to stdout if the verbose flag is set in the format of fmt.Println
func Debugln(msg ...interface{}) {
	if viper.GetBool("always.verbose") {
		debugLog.Println(msg...)
	}
}

// Infof logs to stdout in the format of fmt.Printf
func Infof(format string, v ...interface{}) {
	infoLog.Printf(format, v...)
}

// Warnf logs to stderr in the format of fmt.Printf
func Warnf(format string, v ...interface{}) {
	warnLog.Printf(format, v...)
}

// Errorf logs a string to stderr and exits the program with exit code 1 in the format of fmt.Fatalf
func Errorf(format string, v ...interface{}) {
	err := fmt.Errorf(format, v...)
	errorLog.Fatalln(err)
}

// Debugf logs to stdout if the verbose flag is set in the format of fmt.Printf
func Debugf(format string, v ...interface{}) {
	if viper.GetBool("always.verbose") {
		debugLog.Printf(format, v...)
	}
}
