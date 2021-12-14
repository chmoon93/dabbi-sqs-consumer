package log

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

//
//func Init(sentryDSN, release, logLevel string) {
func Init(release, logLevel string) {
	logrus.SetOutput(os.Stdout)
	logrus.AddHook(&apmlogrus.Hook{})
	SetLogLevel(logLevel)
	//initSentry(sentryDSN, release)
}

//func initSentry(sentryDSN, release string) {
//	if err := sentry.Init(sentry.ClientOptions{
//		Dsn:     sentryDSN,
//		Release: release,
//	}); err != nil {
//		logrus.Errorf("failed to init sentry: %s", err.Error())
//	}
//}

func SetLogLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(msg string) {
	logrus.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Error(msg)
}

func ErrorfWithTracing(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ErrorWithTracing(ctx, msg)
}

func ErrorWithTracing(ctx context.Context, msg string) {
	Error(msg)
}
