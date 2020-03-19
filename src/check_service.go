package main

import (
	"context"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	debug = false
)

func initLogger() {
	var logLevel = logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logrus.AddHook(rotateFileHook)
}

type SupConfig struct {
	Services  []string `yaml:"services"`
	VaultFile string   `yaml:"vaultIniLocation"`
	Disks     []string `yaml:"disks"`
	Port      int      `yaml:"listeningPort"`
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

func checkRunningService(serviceName string, manager *mgr.Mgr) (bool, error) {
	service, err := manager.OpenService(serviceName)
	if err != nil {
		return false, err
	}

	status, err := service.Query()
	if err != nil {
		return false, err
	}

	return status.State == svc.Running, nil
}

func testConnection(addr string, portOptional ...int) error {
	port := 1858
	if len(portOptional) > 0 {
		port = portOptional[0]
	}

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.DialContext(ctx, "tcp", addr+":"+strconv.Itoa(port))

	return err
}

// Find ADDRESS filed in Vault.ini file and return list of IP Address (string format)
func findVaultIPAddress(iniFilePath string) ([]string, error) {
	cfg, err := ini.Load(iniFilePath)
	if err != nil {
		panic(err)
	}

	vaultIPAddresses := cfg.Section("").Key("ADDRESS").String()

	return strings.Split(vaultIPAddresses, ","), nil

}

func findVaultINIFile() ([]string, error) {
	initialSearchList := [...]string{"d:\\Cyberark", "d:\\Program Files\\Cyberark", "d:\\Program Files (x86)\\CyberArk",
		"c:\\Cyberark", "c:\\Program Files\\Cyberark", "c:\\Program Files (x86)\\CyberArk"}

	var finalSearchList []string
	for _, element := range initialSearchList {
		if _, err := os.Stat(element); !os.IsNotExist(err) {
			finalSearchList = append(finalSearchList, element)
			logrus.Debug("keeping " + element)
		}
	}

	var fileList []string
	for _, validPath := range finalSearchList {
		e := filepath.Walk(validPath, func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, "Vault\\Vault.ini") {
				fileList = append(fileList, path)
			}
			return nil
		})

		if e != nil {
			panic(e)
		}
	}

	return fileList, nil
}

func checkServices(services []string) bool {

	// Check running services
	manager, err := mgr.Connect()
	if err != nil {
		logrus.Info("Unable to connect to service manager !")
		return false
		//panic(err)
	}
	defer manager.Disconnect()

	for _, service := range services {
		isRunning, err := checkRunningService(service, manager)
		if err != nil {
			logrus.Fatal("Service " + service + " does not exists or name is invalid ! Exiting.")
		}
		if !isRunning {
			logrus.Info("Service " + service + " is not running! ")
			return false
		}
	}

	return true
}

func getVaultsIPs(vaultIPs []string) []string {
	if len(vaultIPs) == 0 {
		logrus.Info("No IP provided in configuration file, trying to find Vault.ini file")
		list, err := findVaultINIFile()
		if err != nil {
			logrus.Fatal("Unable to find Vault INI file")
		}

		logrus.Info("Using Vault.ini file : " + list[0])

		vaultIPs, _ = findVaultIPAddress(list[0])
	}

	logrus.Info("Using addresses : " + strings.Join(vaultIPs, ","))
	return vaultIPs
}

func checkVault(vaultIPs []string) bool {

	for _, ipAddr := range vaultIPs {
		logrus.Debug("Testing connection to ", ipAddr)
		err := testConnection(ipAddr)
		if err != nil {
			logrus.Debug("Vault connection failed !")
		} else {
			logrus.Debug("Vault connection succeeded !")
			return true
		}
	}

	log.Println("Unable to connect to Vaults !")
	return false
}

type FinalConfig struct {
	services []string
	vaultIPs []string
	disks    []string
	port     int
}

func initialize() FinalConfig {
	var finalConfig FinalConfig
	initLogger()
	supConfig, err := getConf("config.yaml")
	if err != nil {
		logrus.Panic("Error reading configuration file !")
	}

	configVaultsIPs := strings.Split(supConfig.VaultFile, ",")
	finalConfig.vaultIPs = getVaultsIPs(configVaultsIPs)
	finalConfig.services = supConfig.Services
	finalConfig.disks = supConfig.Disks
	finalConfig.port = supConfig.Port

	return finalConfig

}

func status(config FinalConfig) bool {

	logrus.Debug("Starting health tests...")

	vaultConn := checkVault(config.vaultIPs)
	if !vaultConn {
		logrus.Info("Connexion to Vaults failed !")
		return false
	}

	serviceStatus := checkServices(config.services)
	if !serviceStatus {
		logrus.Info("Service check failed !")
		return false
	}

	for _, disk := range config.disks {
		diskStatus := DiskUsage(disk)
		if !diskStatus {
			logrus.Info("Disk " + disk + " available space is less than 10% !")
			return false
		}
	}

	logrus.Debug("All checks went well !")
	return true

}
