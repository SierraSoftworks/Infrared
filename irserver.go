package main

import (
	"github.com/SierraSoftworks/Infrared/lib"
	"github.com/SierraSoftworks/Infrared/lib/config"
	"log"
	"os"
)

func main() {
	configFile := "irserver.json"
	config := config.Server{}
	DefaultConfig(&config)

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration file: %s", err)
	}

	LoadConfigFromCommandLine(&config)

	config.Log()

	VerifyConfig(&config)

	config.Save(configFile)

	server := infrared.Setup(&config)

	server.Start()
}

func DefaultConfig(config *config.Server) {
	if config.Database.Hosts == "" {
		config.Database.Hosts = "localhost"
	}

	if config.Database.Database == "" {
		config.Database.Database = "infrared"
	}
}

func LoadConfigFromCommandLine(config *config.Server) {
	if len(os.Args) > 2 {
		config.ListenOn = os.Args[2]
	}

	if len(os.Args) > 3 {
		config.Database.Hosts = os.Args[3]
	}

	if len(os.Args) > 4 {
		config.Database.Database = os.Args[4]
	}
}

func VerifyConfig(config *config.Server) {
	if config.ListenOn == "" {
		log.Fatal("No valid listening specification provided for the server.")
	}
}
