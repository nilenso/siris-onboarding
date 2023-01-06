package log

import (
	"fmt"
	"os"
)

type StdOutLogger struct {
	level Level
}

func New(level Level) Logger {
	return &StdOutLogger{
		level: level,
	}
}

func (s *StdOutLogger) Log(level Level, message interface{}) {
	if s.IsLevelEnabled(level) {
		fmt.Fprintln(os.Stdout, message)
	}
}

func (s *StdOutLogger) IsLevelEnabled(level Level) bool {
	return level <= s.level
}
