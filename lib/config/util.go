package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func FileExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func SaveJson(filename string, c interface{}) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func LoadJson(filename string, c interface{}) error {
	err := FileExists(filename)
	if err != nil {
		log.Printf("Configuration file %s does not exist", filename)
		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	return nil
}

func LogJson(c interface{}) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	dataBuffer := bytes.NewBuffer(data)
	log.Printf("%s", dataBuffer.String())

	return nil
}
