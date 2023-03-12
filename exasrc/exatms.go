package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var (
	neighbor string
	LogFile  string
	LogLevel string
)

func initLogger() {
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetFormatter(&log.TextFormatter{})
	if LogLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} // else INFO by default
}

func main() {

	flag.StringVar(&LogFile, "logfile", "/var/log/exabgp/exatms.log", "Path to a logfile")
	flag.StringVar(&LogLevel, "loglevel", "INFO", "Logging Level")
	flag.StringVar(&neighbor, "neighbor", "127.0.0.1", "BGP neighbor")

	flag.Parse()

	initLogger()
	log.Info("Starting..")

	go MessageHandler()

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Info("Signal received.. Stopping")

}
