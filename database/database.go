package database

import (
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/marcospsbrito/dic/config"
)

// connect returns an instance of a db connection
func connect(dbURL string, dbName string) (*mgo.Database, error) {
	session, err := mgo.Dial(dbURL)

	if err != nil {
		return nil, err
	}
	return session.DB(dbName), nil
}

// New calls the DB initialization
func New(config config.Config) (*mgo.Database, error) {
	log.WithFields(
		log.Fields{
			"url":      config.MongoURL,
			"database": config.MongoDBName,
		}).Info("opening database connection")

	return connect(config.MongoURL, config.MongoDBName)
}
