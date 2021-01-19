package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	debug = false
)

func findPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func initLogger() {
	var logLevel = logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  findPath() + "/logs/info.log",
		logrus.DebugLevel: findPath() + "/logs/debug.log",
	}

	logrus.SetLevel(logLevel)

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	fmt.Println("Adding hook")
	logrus.AddHook(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{TimestampFormat: time.RFC822},
	))
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.Info("Invalid log level given in configuration file, using info (valid values " +
			"are debug or info")
	}
}

type SupConfig struct {
	Services  []string `yaml:"services"`
	VaultFile string   `yaml:"vaultIniLocation"`
	Disks     []string `yaml:"disks"`
	Port      int      `yaml:"listeningPort"`
	LogLevel  string   `yaml:"logLevel"`
}

// read YAML configuration file
func getConf(fileName string) (SupConfig, error) {
	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Fatal("Unable to read configuration file")
	}
	var supConfig SupConfig

	err = yaml.Unmarshal(confFile, &supConfig)
	if err != nil {
		// probably not valid yaml file
		logrus.Fatal("probably not valid yaml file")
	}

	return supConfig, nil
}

type FinalConfig struct {
	services []string
	vaultIPs []string
	disks    []string
	port     int
}

func initialize() FinalConfig {
	var finalConfig FinalConfig
	//initLogger()
	fmt.Println(findPath())
	supConfig, err := getConf(findPath() + "/config.yaml")

	if err != nil {
		logrus.Panic("Error reading configuration file !")
		//fixing ide linting bug :
		panic(err)
	}

	setLogLevel(supConfig.LogLevel)

	finalConfig.vaultIPs = getVaultsIPs(supConfig.VaultFile)
	finalConfig.services = supConfig.Services
	finalConfig.disks = supConfig.Disks
	finalConfig.port = supConfig.Port

	return finalConfig

}
