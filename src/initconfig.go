package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	var logLevel = log.InfoLevel
	if debug {
		logLevel = log.DebugLevel
	}

	pathMap := lfshook.PathMap{
		log.InfoLevel:  findPath() + "/logs/info.log",
		log.DebugLevel: findPath() + "/logs/debug.log",
	}

	log.SetLevel(logLevel)

	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	fmt.Println("Adding hook")
	log.AddHook(lfshook.NewHook(
		pathMap,
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
	Services  []string `yaml:"services"`
	VaultFile string   `yaml:"vaultIniLocation"`
	Disks     []string `yaml:"disks"`
	Port      int      `yaml:"listeningPort"`
	LogLevel  string   `yaml:"logLevel"`
	VaultsIP  string   `yaml:"vaultsIP"`
}

// read YAML configuration file
func getConf(fileName string) (SupConfig, error) {
	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Unable to read configuration file")
	}
	var supConfig SupConfig

	err = yaml.Unmarshal(confFile, &supConfig)
	if err != nil {
		// probably not valid yaml file
		log.Fatal("probably not valid yaml file")
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
	fmt.Println(findPath())
	supConfig, err := getConf(findPath() + "/config.yaml")

	if err != nil {
		log.Fatal("Error reading configuration file !")
		panic(err)
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

	return finalConfig

}
