package logrushooks

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log/syslog"
	"os"
	"strings"
)

type SyslogHook struct {
	SyslogNetwork string
	SyslogRaddr   string
	Tag           string
	Writer        *syslog.Writer
	Formater      logrus.Formatter
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (syslog.Priority, error) {

	switch strings.ToLower(lvl) {
	case "panic":
		return syslog.LOG_EMERG, nil
	case "fatal":
		return syslog.LOG_CRIT, nil
	case "error":
		return syslog.LOG_ERR, nil
	case "warn", "warning":
		return syslog.LOG_WARNING, nil
	case "info":
		return syslog.LOG_INFO, nil
	case "debug":
		return syslog.LOG_DEBUG, nil
	}

	return 0, fmt.Errorf("invalid log level: '%s' not in [panic, fatal, error, warn|ing, info, debug]", lvl)
}

type Option func(*SyslogHook)

func WithFormater(formater *logrus.Formatter) Option {
	return Option(func(slog *SyslogHook) {
		slog.Formater = formater
	})
}

func WithTag(tag string) Option {
	return Option(func(slog *SyslogHook) {
		slog.Tag = tag
	})
}

/*
Create a hook to be added to an instance of logger. This is called with
   hook, err := NewSyslogHook("udp", "localhost:514", "debug", )
   if err != nil {
       log.Fatalf("Syslog hook init fail: %s", err)
   }
   log.Hooks.Add(hook)
*/
func NewSyslogHook(network, addr, level string, opts ...Option) (*SyslogHook, error) {

	priority, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}

	sLog := &SyslogHook{
		SyslogNetwork: network,
		SyslogRaddr:   addr,
	}

	for _, o := range opts {
		o(sLog)
	}

	if sLog.Formater == nil {
		sLog.Formater = logrus.Logger.Formatter
	}

	w, err := syslog.Dial(network, addr, priority, sLog.Tag)
	if err != nil {
		return nil, err
	}
	sLog.Writer = w

	return sLog, err
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {

	line, err := hook.Formater.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(string(line))
	case logrus.FatalLevel:
		return hook.Writer.Crit(string(line))
	case logrus.ErrorLevel:
		return hook.Writer.Err(string(line))
	case logrus.WarnLevel:
		return hook.Writer.Warning(string(line))
	case logrus.InfoLevel:
		return hook.Writer.Info(string(line))
	case logrus.DebugLevel:
		return hook.Writer.Debug(string(line))
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
