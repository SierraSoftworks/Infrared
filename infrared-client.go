package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var err error
	var res *http.Response

	configServer := os.Getenv("IR_SERVER")
	nodePort, err := strconv.Atoi(os.Getenv("IR_PORT"))

	if err != nil {
		log.Fatalf("IR_PORT could not be deserialized into a valid number: %s", err.Error())
	}

	node := NodeEntry{os.Getenv("IR_NODE"), os.Getenv("IR_NODETYPE"), "", nodePort}

	node.LoadFromConfig()

	if configServer == "" {
		log.Fatal("No IR_SERVER environment variable present.")
	}

	if node.Type == "" {
		log.Fatal("No IR_NODETYPE environment variable specified and no node_type specified in config file.")
	}

	if node.ID == "" {
		log.Printf("Creating new node entry on %s\n", configServer)
		hostname, err := os.Hostname()

		if err != nil {
			log.Fatalf("Failed to determine node's hostname: %s", err.Error())
		}

		node.Hostname = hostname
		creationRequest, err := json.Marshal(NodeCreationRequest{node.Hostname, node.Port})
		if err != nil {
			log.Fatalf("Failed to serialize node creation request: %s", err.Error())
		}

		res, err = http.Post(fmt.Sprintf("%s/api/v1/%s", configServer, node.Type), "application/json", bytes.NewReader(creationRequest))

		if err != nil {
			log.Fatalf("Failed to create new node on %s: %s", configServer, err.Error())
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = json.Unmarshal(buf.Bytes(), &node)

		if err != nil {
			log.Fatalf("Failed to deserialize server response to new node creation: %s", err.Error())
		}

		node.SaveConfig()
	}

	for {
		res, err = http.Get(fmt.Sprintf("%s/api/v1/%s/%s/heartbeat", configServer, node.Type, node.ID))
		if err != nil {
			log.Printf("\nFailed to trigger heartbeat: %s\n", err.Error())
		} else if res.StatusCode != http.StatusOK {
			log.Printf("\nFailed to trigger heartbeat: %s\n", res.Status)
		} else {
			log.Print(".")
		}
		time.Sleep(10 * time.Second)
	}
}

type NodeEntry struct {
	ID       string `json:"id"`
	Type     string `json:"node_type"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func (e NodeEntry) LoadFromConfig() {
	data, err := ioutil.ReadFile("infrared.json")
	if err == nil {
		json.Unmarshal(data, e)
	} else {
		log.Printf("\nFailed to read IR client config file: %s\n", err.Error())
	}
}

func (e NodeEntry) SaveConfig() {
	data, err := json.Marshal(e)
	if err == nil {
		err = ioutil.WriteFile("infrared.json", data, os.ModePerm)
		if err != nil {
			log.Printf("\nFailed to write IR client config file: %s\n", err.Error())
		}
	} else {
		log.Printf("\nFailed to serialize IR client configuration: %s\n", err.Error())
	}
}

type NodeCreationRequest struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}
