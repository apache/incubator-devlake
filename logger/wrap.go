package logger

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/gookit/color.v1"
)

const (
	DEBUG   = "[DEBUG] "
	WARN    = "[WARN] "
	INFO    = "[INFO] "
	SUCCESS = "[SUCCESS] "
	ERROR   = "[ERROR] "
	FATAL   = "[FATAL] "
)

// ROOTCALLER should be main proc
const ROOTCALLER = "main.main"

// REMOVELEVEL logrus stack level is 9, should be removed from stack trace
const REMOVELEVEL = 10

type CliLoggerFormatter struct {
	level           logrus.Level
	showType        string
	formatLevelName string
	prefix          string
}

// Format implement Format interface to output custom log
func (m *CliLoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	m.levelPrintRender()

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	if m.level == logrus.ErrorLevel && logrus.GetLevel() == logrus.DebugLevel {
		entry.Message = addCallStackIgnoreLogrus(entry.Message)
	}

	newLog := fmt.Sprintf("%s %s %s %s\n", timestamp, m.prefix, m.formatLevelName, entry.Message)

	_, err := b.WriteString(newLog)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// levelPrintRender render symbo and level according to type
func (m *CliLoggerFormatter) levelPrintRender() {
	switch m.showType {
	case "debug":
		m.level = logrus.DebugLevel
		m.formatLevelName = color.Blue.Render(DEBUG)
		m.prefix = color.Blue.Render(normal.Debug)
	case "info":
		m.level = logrus.InfoLevel
		m.formatLevelName = color.FgLightBlue.Render(INFO)
		m.prefix = color.FgLightBlue.Render(normal.Info)
	case "warn":
		m.level = logrus.WarnLevel
		m.formatLevelName = color.Yellow.Render(WARN)
		m.prefix = color.Yellow.Render(normal.Warn)
	case "error":
		m.level = logrus.ErrorLevel
		m.formatLevelName = color.BgRed.Render(ERROR)
		m.prefix = color.Red.Render(normal.Error)
	case "fatal":
		m.level = logrus.FatalLevel
		m.formatLevelName = color.BgRed.Render(FATAL)
		m.prefix = color.Red.Render(normal.Fatal)
	case "success":
		m.level = logrus.InfoLevel
		m.formatLevelName = color.Green.Render(SUCCESS)
		m.prefix = color.Green.Render(normal.Success)
	}
}

type SeparatorFormatter struct{}

// Format implement Format interface to output custom log
func (s *SeparatorFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	newLog := fmt.Sprintf("%s %s %s %s\n",
		timestamp,
		color.Blue.Render(normal.Info),
		color.Blue.Render(INFO),
		color.Blue.Render(fmt.Sprintf("%s %s %s", "-------------------- [ ", entry.Message, " ] --------------------")))

	_, err := b.WriteString(newLog)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// addCallStackIgnoreLogrus add call stack to log message without logrus stack
func addCallStackIgnoreLogrus(rawMessage string) string {
	stackMessage := rawMessage
	for i := REMOVELEVEL; ; i++ {
		pc, file, line, _ := runtime.Caller(i)
		stackMessage = stackMessage + "\n  -- " + file + fmt.Sprintf(" %d", line)
		entrance := runtime.FuncForPC(pc).Name()
		if entrance == ROOTCALLER || entrance == "" {
			break
		}
	}
	return stackMessage
}
