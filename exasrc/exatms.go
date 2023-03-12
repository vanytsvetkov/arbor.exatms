package main

import (
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path"
	"runtime"
)

var (
	neighbor   string
	defaultASN string
	LogFile    string
	LogLevel   string
)

func initLogger() {
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
		TrimMessages:  true,
		CallerFirst:   true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			return fmt.Sprintf(" %v::%v", fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function))
		},
	})

	log.SetReportCaller(true)
	if LogLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} // else INFO by default
}

func main() {

	flag.StringVar(&LogFile, "logfile", "/var/log/exabgp/exatms.log", "Path to a logfile")
	flag.StringVar(&LogLevel, "loglevel", "INFO", "Logging Level")
	flag.StringVar(&neighbor, "neighbor", "127.0.0.1", "BGP neighbor")
	flag.StringVar(&defaultASN, "peer_as", "65500", "Default ASN")

	flag.Parse()

	initLogger()
	log.Info("Starting...")

	go MessageHandler()

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Info("Signal received... Stopping")

}
