package logrushooks

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log/syslog"
	"os"
)

type SyslogHook struct {
	Writer     *syslog.Writer
	syslogDial func(network, raddr string, priority syslog.Priority, tag string) (*syslog.Writer, error)

	formater logrus.Formatter
	addr     string
	levels   []logrus.Level
	network  string
	priority syslog.Priority
	tag      string
}

type Option func(*SyslogHook)

func WithFormater(formater logrus.Formatter) Option {
	return func(slog *SyslogHook) {
		slog.formater = formater
	}
}

func WithTag(tag string) Option {
	return func(slog *SyslogHook) {
		slog.tag = tag
	}
}

func WithNetwork(network string) Option {
	return func(slog *SyslogHook) {
		slog.network = network
	}
}

func WithPriority(p syslog.Priority) Option {
	return func(slog *SyslogHook) {
		if p == 0 {
			slog.priority = p
		}
		slog.priority = syslog.LOG_INFO
	}
}

/*
Create a hook to be added to an instance of logger. This is called with
	syslog_host = `127.0.0.1`
    syslog_level = `info
`
    format := (&format.SomeFormater{
		DisableTimestamp:   true,
		MessageAfterFields: true,
	}).Init()

	sysHook, err := logrushooks.NewSyslogHook(
		syslog_host,
		syslog_level,
		logrushooks.WithFormater(format),
	)

	if err != nil {
		logger.Fatalf("Syslog hook init fails: %s", err)
	}

	logger.Hooks.Add(sysHook)
*/
func NewSyslogHook(addr, level string, opts ...Option) (*SyslogHook, error) {

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	levels := []logrus.Level{0}
	for i := 1; i <= int(lvl); i++ {
		levels = append(levels, logrus.Level(i))
	}

	sLog := &SyslogHook{
		addr:       addr,
		levels:     levels,
		priority:   syslog.LOG_INFO,
		network:    "udp",
		formater:   logrus.StandardLogger().Formatter,
		syslogDial: syslog.Dial,
	}

	for _, o := range opts {
		o(sLog)
	}

	w, err := sLog.syslogDial(sLog.network, sLog.addr, sLog.priority, sLog.tag)
	if err != nil {
		return nil, err
	}
	sLog.Writer = w

	return sLog, nil
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	line, err := hook.formater.Format(entry)
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
	return hook.levels
}
