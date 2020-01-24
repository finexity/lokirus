package lokirus

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LokirusHook struct {
	AcceptedLevels []logrus.Level
	Client         promtail.Client
}

func New(hostUrl string, source string) (*LokirusHook, error) {
	labels := "{source=\"" + source + "\"}"
	conf := promtail.ClientConfig{
		PushURL:            hostUrl + "/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.DEBUG,
	}
	client, err := promtail.NewClientProto(conf)

	if err != nil {
		return nil, err
	}
	hook := &LokirusHook{
		AcceptedLevels: logrus.AllLevels,
		Client:         client,
	}

	return hook, nil
}
func (l *LokirusHook) Fire(e *logrus.Entry) error {
	fmt.Printf("LokirusHook2 %v", e.Level)
	line, err := e.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
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
func (l *LokirusHook) Levels() []logrus.Level {
	return l.AcceptedLevels
}
