package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/chmoon93/dabbi-sqs-consumer/config"
	"github.com/chmoon93/dabbi-sqs-consumer/consume"
	"github.com/chmoon93/dabbi-sqs-consumer/log"

	aws "github.com/aws/aws-sdk-go-v2/config"
)

// inject from ldflags
var appName = ""
var buildVersion = ""
var buildCommit = ""
var buildDate = ""

func startConsumeMessagesFromSQS() context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())

	awsConfig, err := aws.LoadDefaultConfig(context.TODO(), aws.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Debug("unable to load SDK config, %v", err)
	}

	go consume.ConsumeMessages(ctx, awsConfig)

	return cancelFunc
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gConfig, err := config.LoadConfig(appName, buildVersion, buildCommit, buildDate)
	log.Debugf("config: %s", err)
	if err != nil {
		fmt.Printf("# Read Config Error: %s\n", err)
		return
	}
	gConfig.PrintConfig()

	log.Init(gConfig.Name, gConfig.LogLevel)

	cancelConsumeMessagesFromSQS := startConsumeMessagesFromSQS()

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Infof("catched signal: %+v", sig)

		done <- true
	}()

	<-done

	cancelConsumeMessagesFromSQS()
	close(quit)
	close(done)

	time.Sleep(1 * time.Second)
	log.Info("# shutdown dabb-sqs-consumer !!!")

}
