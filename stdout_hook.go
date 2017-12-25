package logrushooks

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type StdoutHook struct {
	Writer   io.Writer
	formater logrus.Formatter
	levels   []logrus.Level
}

type StdoutOption func(*StdoutHook)

func WithStdoutFormater(formater logrus.Formatter) StdoutOption {
	return func(slog *StdoutHook) {
		slog.formater = formater
	}
}

func NewStdoutHook(level string, opts ...StdoutOption) (*StdoutHook, error) {

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	levels := []logrus.Level{0}
	for i := 1; i <= int(lvl); i++ {
		levels = append(levels, logrus.Level(i))
	}

	sLog := &StdoutHook{
		levels:   levels,
		formater: logrus.StandardLogger().Formatter,
		Writer:   os.Stderr,
	}

	for _, o := range opts {
		o(sLog)
	}

	return sLog, nil
}

func (hook *StdoutHook) Fire(entry *logrus.Entry) error {
	line, err := hook.formater.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *StdoutHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *StdoutHook) SetFormater(formater logrus.Formatter) {
	hook.formater = formater
}
