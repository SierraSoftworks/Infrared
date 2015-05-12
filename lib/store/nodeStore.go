package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"time"
)

type NodeEntry struct {
	ID       string    `json:"id" bson:"_id"`
	NodeType string    `json:"node_type" bson:"node_type"`
	Hostname string    `json:"hostname" bson:"hostname"`
	Port     int       `json:"port" bson:"port"`
	LastSeen time.Time `json:"lastSeen" bson:"lastSeen"`
}

func (c NodeEntry) valid() bool {
	return len(c.ID) == 16 && len(c.NodeType) > 0 && len(c.Hostname) > 0 && c.Port > 0
}

func (c NodeEntry) FromMap(m bson.M) NodeEntry {
	c.ID = m["_id"].(string)
	c.NodeType = m["node_type"].(string)
	c.LastSeen = m["lastSeen"].(time.Time)
	c.Hostname = m["hostname"].(string)
	c.Port = m["port"].(int)
	return c
}

func (c NodeEntry) ToMap() bson.M {
	return bson.M{"_id": c.ID, "node_type": c.NodeType, "lastSeen": c.LastSeen, "hostname": c.Hostname, "port": c.Port}
}

type NodeCreateRequest struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func CreateNodeEntry(node_type string, request NodeCreateRequest) (NodeEntry, error) {
	id := MakeNodeID(16)

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")

	entry := NodeEntry{id, node_type, request.Hostname, request.Port, time.Now()}
	err = c.Insert(entry.ToMap())

	return entry, err
}

func GetNodeEntry(node_type, id string) (NodeEntry, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")
	node := NodeEntry{}

	err = c.Find(bson.M{"_id": id, "node_type": node_type}).One(&node)

	return node, err
}

func GetNodeEntries(node_type string) ([]NodeEntry, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")

	var nodes []NodeEntry

	err = c.Find(bson.M{"node_type": node_type, "lastSeen": bson.M{"$gt": time.Now().Add(-30 * time.Second)}}).Iter().All(&nodes)

	if nodes == nil {
		nodes = make([]NodeEntry, 0)
	}

	return nodes, err
}

func UpdateNodeEntryLastSeen(node_type, id string) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")
	err = c.Update(bson.M{"_id": id, "node_type": node_type}, bson.M{"$set": bson.M{"lastSeen": time.Now()}})
	return err
}

func UpdateNodeEntry(node_type, id string, request NodeCreateRequest) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")
	err = c.Update(bson.M{"_id": id, "node_type": node_type}, bson.M{"$set": bson.M{"hostname": request.Hostname, "port": request.Port}})
	return err
}

func RemoveNodeEntry(node_type, id string) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	c := session.DB("infrared").C("nodes")

	err = c.Remove(bson.M{"_id": id, "node_type": node_type})
	return err
}

var nodeIDLetters = []rune("0123456789abcdef")

func MakeNodeID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = nodeIDLetters[rand.Intn(len(nodeIDLetters))]
	}
	return string(b)
}
