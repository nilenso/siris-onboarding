package log

import "fmt"

type Level uint32

// Logging levels ordered in decreasing severity
const (
	Fatal Level = iota
	Error
	Warning
	Info
	Debug
)

type Logger interface {
	Log(level Level, message interface{})
}

func (l Level) String() string {
	level, err := levelToString(l)
	if err != nil {
		return "unknown"
	}
	return level
}

func levelToString(level Level) (string, error) {
	switch level {
	case Fatal:
		return "fatal", nil
	case Error:
		return "error", nil
	case Warning:
		return "warning", nil
	case Info:
		return "info", nil
	case Debug:
		return "debug", nil
	}
	return "", fmt.Errorf("Not a valid log level: %v", level)
}
