package lokirus

import (
	"fmt"
	"os"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/sirupsen/logrus"
)

// LokiHook is a logrus hook for Loki
type LokiHook struct {
	AcceptedLevels []logrus.Level
	Client         promtail.Client
}

// New creates a new hook
// hostURL - host of system
// source - the source of the logs
func New(hostURL string, source string) (*LokiHook, error) {
	labels := "{source=\"" + source + "\"}"
	conf := promtail.ClientConfig{
		PushURL:            hostURL + "/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.DEBUG,
	}
	client, err := promtail.NewClientProto(conf)

	if err != nil {
		return nil, err
	}
	hook := &LokiHook{
		AcceptedLevels: logrus.AllLevels,
		Client:         client,
	}

	return hook, nil
}

// Fire forwards the received Logrus entry  to Loki
func (l *LokiHook) Fire(e *logrus.Entry) error {

	// retrieve log message
	line, err := e.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	// execute respective log level
	switch e.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		l.Client.Debugf(line)
	case logrus.InfoLevel:
		l.Client.Infof(line)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		l.Client.Errorf(line)
	default:
		l.Client.Warnf(line)
	}

	return nil
}

// Levels returns all accepted log levels
func (l *LokiHook) Levels() []logrus.Level {
	return l.AcceptedLevels
}
