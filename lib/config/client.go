package config

import (
	"time"
)

type Client struct {
	Server   string    `json:"server"`
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Hostname string    `json:"hostname"`
	Port     int       `json:"port"`
	LastSeen time.Time `json:"lastSeen"`
}

func (c *Client) Save(filename string) error {
	return SaveJson(filename, c)
}

func (c *Client) Load(filename string) error {
	return LoadJson(filename, c)
}

func (c *Client) Log() {
	LogJson(c)
}
