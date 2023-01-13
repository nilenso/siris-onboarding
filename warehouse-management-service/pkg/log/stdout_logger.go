package log

import (
	"fmt"
	"os"
	"time"
)

type StdoutLogger struct {
	level Level
}

func New(level Level) Logger {
	return &StdoutLogger{
		level: level,
	}
}

func (s *StdoutLogger) Log(level Level, message interface{}) {
	if s.IsLevelEnabled(level) {
		fmt.Fprintln(os.Stdout, s.format(level, message))
	}
}

func (s *StdoutLogger) IsLevelEnabled(level Level) bool {
	return level <= s.level
}

func (s *StdoutLogger) format(level Level, message interface{}) string {
	return fmt.Sprintf("%s %v %v", level.String(), time.Now().Format(time.RFC3339), message)
}
