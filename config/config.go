package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Config *WebConfig

const ApplicationJson string = "application/json"
const ContenType string = "Content-Type"

type WebConfig struct {
	WebServerAddr  string
	WebDriverAddr  string
	DriverLogsFile string
}

type ConfigFunc func(*WebConfig)

func DefaultConfig() *WebConfig {
	return &WebConfig{
		WebServerAddr:  "http://localhost:8080",
		WebDriverAddr:  "http://localhost:4444",
		DriverLogsFile: "../driver.logs",
	}
}

func NewConfig(confFn ...ConfigFunc) *WebConfig {
	var conf *WebConfig

	err := loadEnv("../", ".env")
	if err != nil {
		conf = DefaultConfig()

		return conf
	}

	conf = &WebConfig{
		WebServerAddr:  fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		WebDriverAddr:  fmt.Sprintf("%s:%s", os.Getenv("DRIVER_HOST"), os.Getenv("DRIVER_PORT")),
		DriverLogsFile: os.Getenv("DRIVER_LOGS"),
	}

	return conf
}

func WebConfigServerAddr(addr string) ConfigFunc {
	return func(conf *WebConfig) {
		conf.WebServerAddr = addr
	}
}

func WebConfigDriverAddr(addr string) ConfigFunc {
	return func(conf *WebConfig) {
		conf.WebDriverAddr = addr
	}
}

func loadEnv(fRootPath, fName string) error {
	f, err := findFile(fRootPath, fName)
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}

	err = dotenv(f)
	if err != nil {
		return fmt.Errorf("error on dotenv read file: %v", err)
	}

	return nil
}

func dotenv(filepath string) error {
	// read file into memory
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error open file: %v", err)
	}

	defer f.Close()

	// var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		env := scanner.Text()
		if env == "" {
			continue
		}

		key := strings.Split(env, "=")[0]
		value := strings.Split(env, "=")[1]
		os.Setenv(key, value)
	}

	return nil
}

func findFile(fPath, fName string) (string, error) {
	var f string

	err := filepath.WalkDir(fPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println("error on walk dir callback", err)
			return err
		}
		if !info.IsDir() && info.Name() == fName {
			f = path
		}
		return nil
	})

	if err != nil {
		log.Println("error on walk dir", err)
		return "", err
	}

	return f, nil
}
