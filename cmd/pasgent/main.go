package main

import (
	"github.com/Mmx233/Pasgent/pageant"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.StandardLogger()
	logger.Infoln("Pasgent started")
	_ = pageant.Run(logger)
}
