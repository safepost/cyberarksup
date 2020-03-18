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
	debug = true
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
}

// read YAML configuration file
func getConf(fileName string) (SupConfig, error) {
	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Fatal("Unable to read configuration file")
		//panic(err)
	}
	var supConfig SupConfig

	err = yaml.Unmarshal(confFile, &supConfig)
	if err != nil {
		// probably not valid yaml file
		logrus.Fatal("probably not valid yaml file")
		//panic(err)
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

func testVaultConnectivity(addr string, portOptional ...int) error {
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

// Find
func findVaultIPAddress(iniFilePath string) (string, error) {
	cfg, err := ini.Load(iniFilePath)
	if err != nil {
		panic(err)
	}

	return cfg.Section("").Key("ADDRESS").String(), nil

}

func findVaultINIFile() ([]string, error) {
	initialSearchList := [...]string{"d:\\Cyberark", "d:\\Program Files\\Cyberark", "d:\\Program Files (x86)\\CyberArk",
		"c:\\Cyberark", "c:\\Program Files\\Cyberark", "c:\\Program Files (x86)\\CyberArk",
		// for test purpose only !
		"C:\\Users\\gleveill\\Documents\\Bastion\\Splunk"}

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

func checkService(serviceName string) {

	// Check running services
	manager, err := mgr.Connect()
	if err != nil {
		panic(err)
	}
	defer manager.Disconnect()

	isRunning, err := checkRunningService(serviceName, manager)
	if err != nil {
		panic(err)
	}
	if !isRunning {
		logrus.Debug("Service " + serviceName + " is not running! ")
		//Todo : stop listening !
	}
}

func checkVault(vaultIP string) bool {

	if vaultIP == "" {
		// Check Vault connectivity
		list, err := findVaultINIFile()
		if err != nil {
			logrus.Fatal("Unable to find Vault INI file")
		}

		logrus.Debug("Using Vault.ini file : " + list[0])

		//fixme  find vault ini properly
		vaultIP, _ = findVaultIPAddress(list[0])
	}

	vaultIPs := strings.Split(vaultIP, ",")

	for _, ipAddr := range vaultIPs {
		logrus.Debug("Testing connection to ", ipAddr)
		err := testVaultConnectivity(ipAddr)
		if err != nil {
			logrus.Debug("Vault connection failed !")
			//panic(err)
		} else {
			logrus.Debug("Vault connection succeeded !")
			return true
		}
	}

	log.Println("Unable to connect to Vaults !")
	return false
}

func status() bool {
	/*
		path, err := os.Getwd()
		if err != nil {
			log.Fatal("Unable to find working directory")
		}

		file, err := os.OpenFile(path + "/BastionSup.log",os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Unable to create log file !")
		}

		defer file.Close()

		log.SetOutput(file)
	*/

	supConfig, err := getConf("src/config.yaml")
	if err != nil {
		logrus.Fatal("Error reading configuration file !")
		panic(err)
	}

	if supConfig.VaultFile == "" {
		logrus.Info("VaultFile not mentioned in configuration file")
	}

	vaultConn := checkVault(supConfig.VaultFile)

	logrus.Debug(supConfig.VaultFile)
	logrus.Debug(supConfig.Services)

	return vaultConn

}

func main() {
	initLogger()
	logrus.Debug(status())
}
