package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SierraSoftworks/Infrared/lib"
	"github.com/SierraSoftworks/Infrared/lib/config"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)
	config := config.Client{}

	err := config.Load("irclient.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	ParseCommandLine(&config)
	ValidateConfig(&config)

	switch strings.ToLower(os.Args[1]) {
	case "register":
		register(&config)
	case "beam":
		heartbeat(&config)
	}
}

func ParseCommandLine(config *config.Client) {
	if len(os.Args) < 2 {
		log.Fatal("Missing command argument, please specify either 'register' or 'beam'.")
	}

	if len(os.Args) > 2 {
		config.Server = os.Args[2]
	}

	if len(os.Args) > 3 {
		config.Type = os.Args[3]
	}

	if len(os.Args) > 4 {
		config.Hostname = os.Args[4]
	}

	if len(os.Args) > 5 {
		port, err := strconv.Atoi(os.Args[5])
		if err != nil {
			log.Fatalf("Port '%s' could not be converted to a valid number", os.Args[5])
		}
		config.Port = port
	}

	if config.Hostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Failed to determine hostname of current system. Please specify one on the command line or in the irclient.json file.\n%s", err)
		}
		config.Hostname = hostname
	}
}

func ValidateConfig(config *config.Client) {
	if config.Server == "" {
		log.Fatal("Missing server URL, please specify the server's URL on the command line or in the irclient.json file.")
	}

	if config.Type == "" {
		log.Fatal("Missing node type, please specify the node's type on the command line or in the irclient.json file.")
	}

	if config.Port < 1 {
		log.Fatal("Missing service port, please specify the service's port on the command line or in the irclient.json file.")
	}
}

func register(config *config.Client) {
	creationRequest, err := json.Marshal(NodeCreationRequest{config.Hostname, config.Port})
	if err != nil {
		log.Fatalf("Failed to serialize node creation request: %s", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/api/v1/%s", config.Server, config.Type), "application/json", bytes.NewReader(creationRequest))

	if err != nil {
		log.Fatalf("Failed to create new node on %s: %s", config.Server, err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Failed to create new node on %s: %s", config.Server, res.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	err = json.Unmarshal(buf.Bytes(), config)

	if err != nil {
		log.Fatalf("Failed to deserialize server response to new node creation: %s", err)
	}

	err = config.Save("irclient.json")
	if err != nil {
		log.Fatalf("Failed to store node configuration in irclient.json: %s", err)
	}

	log.Printf("Node registered on %s with ID:%s", config.Server, config.ID)
}

func heartbeat(config *config.Client) {
	if config.ID == "" {
		log.Fatal("This client does not yet have an ID associated with it, please run irclient register first.")
	}

	serverUrl, err := url.Parse(config.Server)
	if err != nil {
		log.Fatalf("Failed to parse remote server URL: %s", err)
	}

	ServerAddr, err := net.ResolveUDPAddr("udp", serverUrl.Host)
	if err != nil {
		log.Fatalf("Failed to resolve remote server address for %s: %s", config.Server, err)
	}

	log.Printf("Resolved remote server %s to %s", config.Server, ServerAddr.String())

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Failed to resolve local address: %s", err)
	}

	heartbeat := infrared.Heartbeat{}
	heartbeat.Id = &config.ID
	heartbeat.NodeType = &config.Type
	heartbeatData, err := proto.Marshal(&heartbeat)

	for {
		time.Sleep(1 * time.Second)
		conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)

		if err != nil {
			continue
		}

		_, err = conn.Write(heartbeatData)
		if err != nil {
			log.Printf("Failed to trigger heartbeat: %s", err)
		} else {
			log.Printf("Heartbeat sent")
		}
	}
}

type NodeCreationRequest struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}
