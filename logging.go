package tracelog

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	NormalLogLevel = "NORMAL"
	DevelLogLevel  = "DEVEL"
	timeFlags      = log.LstdFlags | log.Lmicroseconds

	infoPrefix    = "INFO: "
	warningPrefix = "WARNING: "
	errorPrefix   = "ERROR: "
	debugPrefix   = "DEBUG: "
)

var InfoLogger = NewErrorLogger(os.Stdout, infoPrefix)
var WarningLogger = NewErrorLogger(os.Stdout, warningPrefix)
var ErrorLogger = NewErrorLogger(os.Stderr, errorPrefix)
var DebugLogger = NewErrorLogger(ioutil.Discard, debugPrefix)

var LogLevels = []string{NormalLogLevel, DevelLogLevel}
var logLevel = NormalLogLevel
var logLevelFormatters = map[string]string{
	NormalLogLevel: "%v",
	DevelLogLevel:  "%+v",
}

func setupLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewErrorLogger(ioutil.Discard, debugPrefix)
	} else {
		DebugLogger = NewErrorLogger(os.Stdout, debugPrefix)
	}
}

type LogLevelError struct {
	error
}

func NewLogLevelError() LogLevelError {
	return LogLevelError{errors.Errorf("got incorrect log level: '%s', expected one of: '%v'", logLevel, LogLevels)}
}

func (err LogLevelError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

func GetErrorFormatter() string {
	return logLevelFormatters[logLevel]
}

func UpdateLogLevel(newLevel string) error {
	isCorrect := false
	for _, level := range LogLevels {
		if newLevel == level {
			isCorrect = true
		}
	}
	if !isCorrect {
		return NewLogLevelError()
	}

	logLevel = newLevel
	setupLoggers()
	return nil
}

func RedirectLogging(infoWriter, warningWriter, errorWriter, debugWriter io.Writer) {
	InfoLogger = NewErrorLogger(infoWriter, infoPrefix)
	WarningLogger = NewErrorLogger(warningWriter, warningPrefix)
	ErrorLogger = NewErrorLogger(errorWriter, errorPrefix)
	DebugLogger = NewErrorLogger(debugWriter, debugPrefix)
}
