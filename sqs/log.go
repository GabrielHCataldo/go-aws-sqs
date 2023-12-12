package sqs

import "github.com/GabrielHCataldo/go-logger/logger"

func loggerInfo(debugMode bool, v ...any) {
	if debugMode {
		logger.InfoSkipCaller(2, v...)
	}
}

func loggerErr(debugMode bool, v ...any) {
	if debugMode {
		logger.ErrorSkipCaller(2, v...)
	}
}
