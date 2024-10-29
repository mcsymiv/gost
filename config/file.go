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

func loadEnv(fRootPath, fName string) error {
	f, err := FindFile(fRootPath, fName)
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
		if env == "" || strings.ContainsRune(env, '#') { // # serves as comment token, ignored by dotenv
			continue
		}

		key := strings.Split(env, "=")[0]
		value := strings.Split(env, "=")[1]
		os.Setenv(key, value)
	}

	return nil
}

func FindFile(fPath, fName string) (string, error) {
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

func GetPath(dirName string) string {
	return fmt.Sprintf("%s/%s", Root, dirName)
}

func GetRoot(dirName string) string {
	return fmt.Sprintf("%s/", Root)
}
