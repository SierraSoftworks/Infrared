package store

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type ConfigEntry struct {
	NodeType      string      `json:"node_type" bson:"_id"`
	LastUpdated   time.Time   `json:"lastUpdated" bson:"lastUpdated"`
	Configuration interface{} `json:"config" bson:"config"`
}

func (c ConfigEntry) valid() bool {
	return len(c.NodeType) > 0
}

func (c ConfigEntry) FromMap(m bson.M) ConfigEntry {
	c.NodeType = m["_id"].(string)
	c.LastUpdated = m["lastUpdated"].(time.Time)
	c.Configuration = m["config"]
	return c
}

func (c ConfigEntry) ToMap() bson.M {
	return bson.M{"_id": c.NodeType, "lastUpdated": c.LastUpdated, "config": c.Configuration}
}

func GetConfigEntry(node_type string) (interface{}, error) {
	if node_type == "" {
		return nil, ValidationError{"bad request"}
	}

	session, c, err := GetSession("config")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	config := ConfigEntry{}

	err = c.Find(bson.M{"_id": node_type}).One(&config)

	return config.Configuration, err
}

func SetConfigEntry(node_type string, config interface{}) error {
	if node_type == "" {
		return ValidationError{"bad request"}
	}

	session, c, err := GetSession("config")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	_, err = c.Upsert(bson.M{"_id": node_type}, bson.M{"$set": bson.M{"config": config, "lastUpdated": time.Now()}})
	return err
}
