package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type configHandler struct {
	*appConfig
}

var gConfig *configHandler

func init() {
	gConfig = &configHandler{}
}

type appConfig struct {
	// App Property
	Name         string // NAME
	BuildVersion string // Build Version
	BuildCommit  string // Build Commit
	BuildDate    string // Build Date
	LogLevel     string // LOG_LEVEL

	// AWS SQS
	SQSName     string
	SQSWaitTime uint

	AccessKey       string
	SecretAccessKey string

	// Database
	// DbIP       string //	MARIA_DB_IP
	// DbPort     uint   //	MARIA_DB_PORT
	// DbName     string //	MARIA_DB_NAME
	// DbUserName string //	MARIA_DB_USERNAME
	// DbPassword string //	MARIA_DB_PASSWORD

	// // Kafka
	// KafkaTopic   string   // KAFKA_ALERT_TOPIC
	// KafkaBrokers []string // KAFKA_BROKERS

	// CLOSE DURATION
	// AutoCloseDuration uint // EVENT_ALERT_AUTO_CLOSE_DURATION
	// AlertCheckYn      string
}

func (config *configHandler) setConfigByDefault(appName string, buildVersion string, buildCommit string, buildDate string) {
	config.appConfig = &appConfig{
		Name:         appName,
		BuildVersion: buildVersion,
		BuildCommit:  buildCommit,
		BuildDate:    buildDate,
		LogLevel:     "DEBUG",
		SQSName:      "",
		SQSWaitTime:  0,
	}
}

func (config *configHandler) PrintConfig() {
	fmt.Printf("+-------------------------------------+\n")
	fmt.Printf("# Name             : %s\n", gConfig.Name)
	fmt.Printf("# Build Version    : %s\n", gConfig.BuildVersion)
	fmt.Printf("# Build Commit     : %s\n", gConfig.BuildCommit)
	fmt.Printf("# Build Date       : %s\n", gConfig.BuildDate)
	fmt.Printf("# LogLevel         : %s\n", gConfig.LogLevel)

	fmt.Printf("# SQSName          : %s\n", gConfig.SQSName)
	fmt.Printf("# SQSWaitTime      : %d\n", gConfig.SQSWaitTime)

	fmt.Printf("# AccessKey        : %s\n", gConfig.AccessKey)
	fmt.Printf("# SecretAccessKey  : %s\n", gConfig.SecretAccessKey)

	// fmt.Printf("# DbIp             : %s\n", gConfig.DbIP)
	// fmt.Printf("# DbPort           : %d\n", gConfig.DbPort)
	// fmt.Printf("# DbName           : %s\n", gConfig.DbName)
	// fmt.Printf("# DbUserName       : %s\n", gConfig.DbUserName)
	// // fmt.Printf("# DbPassword       : %s\n", gConfig.DbPassword)

	// fmt.Printf("# KafkaTopic       : %s\n", gConfig.KafkaTopic)
	// fmt.Printf("# AutoCloseDuration: %d(sec)\n", gConfig.AutoCloseDuration)
	// fmt.Printf("# AlertCheckYn     : %s\n", gConfig.AlertCheckYn)
	fmt.Printf("+-------------------------------------+\n")
}

func GetConfigHandler() *configHandler {
	return gConfig
}

func LoadConfig(appName string, buildVersion string, buildCommit string, buildDate string) (*configHandler, error) {
	if len(appName) == 0 {
		if err := loadEnvAsStr(&appName, "NAME", REQUIRE_TRUE); err != nil {
			return nil, err
		}
	}

	gConfig.setConfigByDefault(appName, buildVersion, buildCommit, buildDate)

	// LOG_LEVEL
	if err := loadEnvAsStr(&gConfig.LogLevel, "LOG_LEVEL", REQUIRE_FALSE); err != nil {
		return nil, err
	}

	// SQS_NAME
	if err := loadEnvAsStr(&gConfig.SQSName, "SQS_NAME", REQUIRE_TRUE); err != nil {
		return nil, err
	}

	// SQS_WAIT_TIME
	if err := loadEnvAsUInt(&gConfig.SQSWaitTime, "SQS_WAIT_TIME", REQUIRE_TRUE); err != nil {
		return nil, err
	}

	// AWS_ACCESS_KEY_ID
	if err := loadEnvAsStr(&gConfig.AccessKey, "AWS_ACCESS_KEY_ID", REQUIRE_FALSE); err != nil {
		return nil, err
	}

	// AWS_SECRET_ACCESS_KEY
	if err := loadEnvAsStr(&gConfig.SecretAccessKey, "AWS_SECRET_ACCESS_KEY", REQUIRE_FALSE); err != nil {
		return nil, err
	}

	return gConfig, nil
}

func loadEnvAsStr(configVal *string, envKey string, isRequired bool) error {
	envVal := os.Getenv(envKey)
	if envVal == "" && isRequired {
		return fmt.Errorf("%s required", envKey)
	}

	if envVal != "" {
		*configVal = envVal
	}

	return nil
}

func loadEnvAsStrSlice(configVal *[]string, envKey string, isRequired bool) error {
	envVal := os.Getenv(envKey)
	if envVal == "" && isRequired {
		return fmt.Errorf("%s required", envKey)
	}

	if envVal != "" {
		*configVal = strings.Split(envVal, ",")
	}

	return nil
}

func loadEnvAsUInt(configVal *uint, envKey string, isRequired bool) error {
	envVal, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil && isRequired {
		return fmt.Errorf(" %s required", envKey)
	}

	if err == nil {
		*configVal = uint(envVal)
	}

	return nil
}

func loadIntEnv(configVal *uint, envKey string, defaultValue uint) {
	*configVal = defaultValue
	envVal, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		fmt.Printf(" %s required", envKey)
	} else {
		*configVal = uint(envVal)
	}
}
