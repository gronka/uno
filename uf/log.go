package uf

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

var loggingLevel int = 0
var isProd bool = false

func FlashError(err error) {
	if err != nil {
		Error(err)
	}
}

func LoggingError(i interface{}) {
	//TODO: how do we log/notify on a logging error?
	msg := fmt.Sprintf("%v", i)
	logMailer.sendLog("LOGFAILURE!!!", msg)
}

//TODO: pass interface as a list: ...interface{} to add JSON

func Panic(i interface{}) { // 70
	msg := fmt.Sprintf("%v", i)
	log.Panic().Msg(msg)
	logMailer.sendLog("PANIC", msg)
	panic(i)
}

func Fatal(i interface{}) { // 60
	msg := fmt.Sprintf("%v", i)
	log.Fatal().Msg(msg)
	logMailer.sendLog("FATAL", msg)
}

func Error(i interface{}) { // 50
	msg := fmt.Sprintf("%v", i)
	log.Error().Msg(msg)
	logMailer.sendLog("ERROR", msg)
}

func Warn(i interface{}) { // 40
	if loggingLevel < 40 {
		msg := fmt.Sprintf("%v", i)
		log.Warn().Msg(msg)
		logMailer.sendLog("WARN", msg)
	}
}

func Info(i interface{}) { // 30
	if loggingLevel < 30 {
		log.Info().Msg(fmt.Sprintf("%v", i))
	}
}

func Debug(i interface{}) { // 20
	if loggingLevel < 20 {
		log.Debug().Msg(fmt.Sprintf("%v", i))
	}
}

func Trace(i interface{}) { // 10
	if loggingLevel < 10 {
		log.Trace().Msg(fmt.Sprintf("%v", i))
	}
}

var LevelPanic int = 70
var LevelFatal int = 60
var LevelError int = 50
var LevelWarn int = 40
var LevelInfo int = 30
var LevelDebug int = 20
var LevelTrace int = 10

func Glog(gibs *Gibs, glog GlogStruct) { // 10
	glog.ToddId = gibs.ToddId
	glog.UfId = gibs.UfId
	glog.Service = gibs.Conf.UfName
	glog.Time = NowStamp()

	if glog.Level == 0 {
		Fatal(glog.Code)
		Fatal(glog.Msg)
		Fatal("loglevel for a log is zero!")
	}

	if glog.Interface != nil {
		if glog.Msg == "" {
			glog.Msg = fmt.Sprintf("%v", glog.Interface)
		} else {
			glog.Msg = fmt.Sprintf(glog.Msg+";\n%v", glog.Interface)
		}
	}

	ures := GlogCreate(gibs, glog)

	if ures.Errored() {
		LoggingError(ures.Errors)
	}

	// NOTE: log locally last in case we call panic()
	switch glog.Level {
	case LevelPanic: // 70
		//NOTE: do a LoggingError to help ensure we get a log somewhere
		LoggingError(glog.Interface)
		Panic(glog.Interface)

	case LevelFatal: // 60
		LoggingError(glog.Interface)
		Fatal(glog.Interface)

	case LevelError: // 50
		Error(glog.Interface)

	case LevelWarn: // 40
		Warn(glog.Interface)

	case LevelInfo: // 30
		Info(glog.Interface)

	case LevelDebug: // 20
		Debug(glog.Interface)

	case LevelTrace: // 10
		Trace(glog.Interface)
	}

}
