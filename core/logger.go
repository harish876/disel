package disel

import (
	"io"
	"log"
	"os"
)

var (
	DEFAULT_WRITER = os.Stdout
	DEFAULT_PREFIX = "Disel-Logger: "
	DEFAULT_FLAGS  = log.Ldate | log.Ltime
	DEFAULT_LEVEL  = 2
)

var (
	TRACE = 0
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
)

type Logger struct {
	writer io.Writer
	prefix string
	flag   int
	logger *log.Logger
	level  int
}

func InitLogger() *Logger {
	return &Logger{
		writer: DEFAULT_WRITER,
		prefix: DEFAULT_PREFIX,
		flag:   DEFAULT_FLAGS,
		logger: nil,
		level:  DEFAULT_LEVEL,
	}
}

func (l *Logger) SetWriter(writer io.Writer) *Logger {
	if writer != DEFAULT_WRITER {
		l.writer = writer
	}
	return l
}

func (l *Logger) SetPrefix(prefix string) *Logger {
	if prefix != DEFAULT_PREFIX {
		l.prefix = prefix
	}
	return l
}

func (l *Logger) SetFlags(flag int) *Logger {
	if flag != DEFAULT_FLAGS {
		l.flag = flag
	}
	return l
}

func (l *Logger) SetLevel(level int) *Logger {
	if level != DEFAULT_LEVEL {
		l.level = level
	}
	return l
}

func (l *Logger) Build() *Logger {
	l.logger = log.New(l.writer, l.prefix, l.flag)
	return l
}

func (l *Logger) Trace(args ...any) {
	if l.level <= TRACE {
		l.logger.Println(args...)
	}
}

func (l *Logger) Debug(args ...any) {
	if l.level <= DEBUG {
		l.logger.Println(args...)
	}
}

func (l *Logger) Info(args ...any) {
	if l.level <= INFO {
		l.logger.Println(args...)
	}
}

func (l *Logger) Warn(args ...any) {
	if l.level <= WARN {
		l.logger.Println(args...)
	}
}

func (l *Logger) Error(args ...any) {
	if l.level <= ERROR {
		l.logger.Println(args...)
	}
}

func (l *Logger) Fatal(args ...any) {
	if l.level <= FATAL {
		l.logger.Println(args...)
	}
}

func (l *Logger) Tracef(format string, args ...any) {
	if l.level <= TRACE {
		l.logger.Printf(format, args...)
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	if l.level <= DEBUG {
		l.logger.Printf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...any) {
	if l.level <= INFO {
		l.logger.Printf(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...any) {
	if l.level <= WARN {
		l.logger.Printf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.level <= ERROR {
		l.logger.Printf(format, args...)
	}
}

func (l *Logger) Fatalf(format string, args ...any) {
	if l.level <= FATAL {
		l.logger.Fatalf(format, args...)
	}
}
