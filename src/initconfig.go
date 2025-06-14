package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigFile = "config.yaml"
	DefaultLogsDir    = "logs"
	DefaultInfoLog    = "CyberarkSupervision.log"
	DefaultDebugLog   = "CyberarkSupervision_debug.log"
)

// If the DEBUG environment variable is set to true, we will not perform checks
var debug = os.Getenv("DEBUG") == "true"

func findPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func initLogger() {
	var logLevel = log.InfoLevel
	if debug {
		logLevel = log.DebugLevel
	}

	// Create the log folder if it does not exist
	logsDir := filepath.Join(findPath(), DefaultLogsDir)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal("Unable to create the log folder if it does not exist: ", err)
	}

	infoPathMap := lfshook.PathMap{
		log.InfoLevel: filepath.Join(logsDir, DefaultInfoLog),
	}

	// PathMap supplémentaire pour écrire info ET debug dans debug.log
	debugPathMap := lfshook.PathMap{
		log.InfoLevel:  filepath.Join(logsDir, DefaultDebugLog),
		log.DebugLevel: filepath.Join(logsDir, DefaultDebugLog),
	}

	log.SetLevel(logLevel)

	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	log.AddHook(lfshook.NewHook(
		infoPathMap,
		&log.JSONFormatter{TimestampFormat: time.RFC822},
	))

	log.AddHook(lfshook.NewHook(
		debugPathMap,
		&log.JSONFormatter{TimestampFormat: time.RFC822},
	))

}

func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.Info("Invalid log level given in configuration file, using info (valid values " +
			"are debug or info")
	}
}

type SupConfig struct {
	Services            []string `yaml:"services"`
	VaultFile           string   `yaml:"vaultIniLocation"`
	Disks               []string `yaml:"disks"`
	Port                int      `yaml:"listeningPort"`
	HealthCheckInterval int      `yaml:"healthCheckInterval"`
	ConnectionTimeout   int      `yaml:"connectionTimeout"`
	LogLevel            string   `yaml:"logLevel"`
	VaultsIP            string   `yaml:"vaultsIP"`
}

// read the YAML configuration file
func getConf(fileName string) (SupConfig, error) {
	var supConfig SupConfig

	confFile, err := os.ReadFile(fileName)
	if err != nil {
		return supConfig, fmt.Errorf("unable to read configuration file %s: %w", fileName, err)
	}

	if err := yaml.Unmarshal(confFile, &supConfig); err != nil {
		return supConfig, fmt.Errorf("invalid YAML format: %w", err)
	}

	return supConfig, nil
}

type FinalConfig struct {
	services            []string
	vaultIPs            []string
	disks               []string
	port                int
	HealthCheckInterval int
	ConnectionTimeout   int
}

func initialize() FinalConfig {
	var finalConfig FinalConfig
	fmt.Println(findPath())
	supConfig, err := getConf(findPath() + "/" + DefaultConfigFile)

	if err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	// Validate configuration
	if supConfig.Port <= 0 || supConfig.Port > 65535 {
		log.Fatal("Invalid port number in configuration")
	}

	if supConfig.HealthCheckInterval <= 1 || supConfig.HealthCheckInterval > 500 {
		log.Fatal("healthCheckInterval should be an integer between 1 and 500 (seconds)")
	}

	if supConfig.ConnectionTimeout <= 1 || supConfig.ConnectionTimeout > 50 {
		log.Fatal("connectionTimeout should be an integer between 1 and 50 (seconds)")
	}

	if len(supConfig.Services) == 0 {
		log.Warn("No services configured for monitoring")
	}

	setLogLevel(supConfig.LogLevel)
	if supConfig.VaultsIP != "" {
		log.Debug("Vault IP were given in config file, using it :" + supConfig.VaultsIP)
		finalConfig.vaultIPs = strings.Split(supConfig.VaultsIP, ",")
	} else {
		log.Debug("Vault IP were NOT given in config file")
		finalConfig.vaultIPs = getVaultsIPs(supConfig.VaultFile)
	}

	finalConfig.services = supConfig.Services
	finalConfig.disks = supConfig.Disks
	finalConfig.port = supConfig.Port
	finalConfig.HealthCheckInterval = supConfig.HealthCheckInterval
	finalConfig.ConnectionTimeout = supConfig.ConnectionTimeout

	return finalConfig

}

func (c FinalConfig) getHealthCheckDuration() time.Duration {
	return time.Duration(c.HealthCheckInterval) * time.Second
}

func (c FinalConfig) getVaultTimeOutDuration() time.Duration {
	return time.Duration(c.ConnectionTimeout) * time.Second
}
