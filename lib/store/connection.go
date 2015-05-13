package store

import (
	"github.com/SierraSoftworks/Infrared/lib/config"
	"gopkg.in/mgo.v2"
)

var configuration *config.Server

func GetSession(collection string) (*mgo.Session, *mgo.Collection, error) {
	session, err := mgo.Dial(configuration.Database.Hosts)
	if err != nil {
		return session, nil, err
	}

	c := session.DB(configuration.Database.Database).C(collection)

	return session, c, err
}

func SetConfig(c *config.Server) {
	configuration = c
}
