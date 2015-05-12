package config

import (
	"time"
)

type Client struct {
	Server string `json:"server"`
	ID string `json:"id"`
	Type string `json:"type"`
	Hostname string `json:"hostname"`
	Port string `json:"port"`
	LastSeen time.Time `json:"lastSeen"`
}

func (c Client) Save(filename string) error {
	return Save(filename, c)
}

func (c Client) Load(filename string) error {
	return Load(filename, c)
}