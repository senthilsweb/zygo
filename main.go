package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/senthilsweb/zygo/pkg/router"
	"github.com/senthilsweb/zygo/pkg/utils"

	"github.com/apex/gateway"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

var (
	flagPort int
	flagEnv  string
)

// staticFileServer implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type staticFileServer struct {
	contents      embed.FS
	staticDirPath string
	indexFileName string
}

//go:embed templates/*
var templateContents embed.FS

func init() {
	log.Info("Application init function start")

	log.Info("Initialize Logger")
	initLogger()
	log.Info("Initialized Logger")

	log.Info("Initialize Configuration")
	//config.Setup()
	log.Info("Initialized Configuration")

	log.Info("Initialize command line args")
	flag.IntVar(&flagPort, "p", -1, "port number for the api server")
	flag.StringVar(&flagEnv, "e", "dev", "Development or Production")

	log.Info("Port = [" + fmt.Sprintf(":%d", flagPort) + "]")
	log.Info("Environment = [" + flagEnv + "]")

	log.Info("Initialized command line args")

	log.Info("Application init function end.")
}

func main() {

	flag.Parse()

	//go func() {
	//	go task.SubscribeAndReceiveMessage()
	//}()

	startServer()

}

func initLogger() {

	logfilepath := utils.AppExecutionPath() + "/" + os.Args[0] + ".log"
	log.Info("logfilepath = " + logfilepath)
	// Set the Lumberjack logger
	ljack := &lumberjack.Logger{
		Filename:   logfilepath,
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     3,
		LocalTime:  true,
	}
	log.Info(ljack)
	//log := logrus.New()
	//
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	//mWriter := io.MultiWriter(os.Stdout, ljack)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"app":             os.Args[0],
		"Runtime Version": runtime.Version(),
		"Number of CPUs":  runtime.NumCPU(),
		"Arch":            runtime.GOARCH,
	}).Info("Application Initializing")
}

func startServer() {
	r := router.Setup()
	//log.Fatal(r)
	//log.Fatal(http.ListenAndServe(":"+port, a.Negroni))
	listener := gateway.ListenAndServe
	portStr := "n/a"
	if flagPort != -1 {
		portStr = fmt.Sprintf(":%d", flagPort)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./public")))
	}
	log.Info("Starting the server " + strconv.Itoa(flagPort))
	done := make(chan bool)
	//go listener(flagPort, r)
	go listener(portStr, r)
	log.Info("Server started at port " + strconv.Itoa(flagPort))
	<-done
}
