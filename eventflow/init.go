package eventflow

import (
	"github.com/c3sr/config"
	"github.com/c3sr/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry = logger.New().WithField("pkg", "evaluation/eventflow")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "evaluation/eventflow")
	})
}
