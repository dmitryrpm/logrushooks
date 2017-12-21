package logrushooks

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"log/syslog"
	"testing"
)

type Deal func(network, raddr string, priority syslog.Priority, tag string) (*syslog.Writer, error)

type MockDeal struct {
	buff *bytes.Buffer
}

func (m MockDeal) Deal(network, raddr string, priority syslog.Priority, tag string) (*syslog.Writer, error) {
	m.buff = bytes.NewBuffer([]byte{})
	return &syslog.Writer{}, nil
}

func WithSyslogDeal(mockDeal Deal) Option {
	return Option(func(slog *SyslogHook) {
		slog.syslogDial = mockDeal
	})
}

func TestCorrectHookWithFire(t *testing.T) {
	mockDeal := MockDeal{}
	sysHook, err := NewSyslogHook(
		`127.0.0.1`,
		`info`,
		WithSyslogDeal(mockDeal.Deal),
	)

	if err != nil || sysHook == nil {
		t.Fatalf("Syslog hook init fails: %s", err)
	}

	entry := logrus.Entry{
		Message: "test",
	}

	err = sysHook.Fire(&entry)
	if err != nil {
		t.Fatalf("Syslog incorrect call Fire: %s", err)
	}
}

func TestIncorrectLogLevelHook(t *testing.T) {
	mockDeal := MockDeal{}
	sysHook, err := NewSyslogHook(
		`127.0.0.1`,
		`trololo`,
		WithSyslogDeal(mockDeal.Deal),
	)

	if err == nil || sysHook != nil {
		t.Fatal("Incorrect create SyslogHook with LogLevel")
	}
}
