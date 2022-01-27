package pkg

import (
    "log"
)

const (
    infoLevel    = "INFO"
    debugLevel   = "DEBUG"
    warningLevel = "WARNING"
    errorLevel   = "Error"
)

var (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorWhite  = "\033[97m"
    // colorGreen  = "\033[32m"
    // colorPurple = "\033[35m"
    // colorCyan   = "\033[36m"
    // colorGray   = "\033[37m"


    logLevelColorMapping = map[string]string{
        infoLevel:    colorWhite,
        debugLevel:   colorBlue,
        warningLevel: colorYellow,
        errorLevel:   colorRed,
    }
)

// Error describes the error type for the tool.
type Error struct {
    Err  error
}

// LogInfo logs the respective error on 'info' level.
func (e *Error) LogInfo() {
    e.log(infoLevel)
}

// LogDebug logs the respective error on 'debug' level.
func (e *Error) LogDebug() {
    e.log(debugLevel)
}

// LogWarning logs the respective error on 'warning' level.
func (e *Error) LogWarning() {
    e.log(warningLevel)
}

// LogError logs the respective error on 'error' level.
func (e *Error) LogError() {
    e.log(errorLevel)
}

// log the respective error structure.
func (e *Error) log(logLevel string) {
    log.Printf("%v[%v] Error message: %v%v\n", logLevelColorMapping[logLevel], logLevel, e.Err, colorReset)
}

// LogMessage describes the log message to print out.
type LogMessage string

// LogInfo logs the respective error on 'info' level.
func (l *LogMessage) LogInfo() {
    l.log(infoLevel)
}

// LogDebug logs the respective error on 'debug' level.
func (l *LogMessage) LogDebug() {
    l.log(debugLevel)
}

// LogWarning logs the respective error on 'warning' level.
func (l *LogMessage) LogWarning() {
    l.log(warningLevel)
}

// LogError logs the respective error on 'error' level.
func (l *LogMessage) LogError() {
    l.log(errorLevel)
}

// log the respective error structure.
func (l *LogMessage) log(logLevel string) {
    log.Printf("%v[%v] Error message: %v%v\n", logLevelColorMapping[logLevel], logLevel, l, colorReset)
}
