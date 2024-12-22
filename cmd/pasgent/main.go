package main

import (
	"github.com/Mmx233/Pasgent/pageant"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.Infoln("pasgent started")
	_ = pageant.Run(logger)
}
