package gored

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SierraSoftworks/Infrared/lib/protocol"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type NodeRegistration struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type NodeDetails struct {
	ID       string    `json:"id"`
	Type     string    `json:"node_type"`
	Hostname string    `json:"hostname"`
	Port     int       `json:"port"`
	LastSeen time.Time `json:"lastSeen"`
}

type APIError io.Reader

func (e APIError) Error() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(io.Reader(e))

	var data map[string]interface{}
	json.Unmarshal(buf, data)
	return fmt.Sprintf("%s: %s", data["error"], data["message"])
}

func Register(server, node_type string, details NodeRegistration) (NodeDetails, error) {
	data, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(fmt.Sprintf("%s/api/v1/%s", server, node_type), "application/json", data)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, APIError(response.Body)
	}

	result := NodeDetails{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	err = json.Unmarshal(buf, &result)

	return result, err
}

func Heartbeat(server, node_type, node_id string) error {
	serverUrl, err := url.Parse(server)
	if err != nil {
		return err
	}
	ServerAddr, err := net.ResolveUDPAddr("udp", serverUrl.Host)
	if err != nil {
		return err
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	heartbeat := protocol.Heartbeat{}
	heartbeat.Id = &node_id
	heartbeat.NodeType = &node_type
	heartbeatData, err := proto.Marshal(&heartbeat)

	conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		return err
	}
	_, err = conn.Write(heartbeatData)
	if err != nil {
		return err
	}

	return nil
}
