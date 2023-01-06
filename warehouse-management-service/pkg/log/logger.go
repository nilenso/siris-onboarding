package log

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
