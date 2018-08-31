package repository

import (
	"crypto/tls"
	"github.com/BetaProjectWave/kube-vault-plugin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tag-service/logger"
	"github.com/tag-service/model"
	cfg "github.com/tag-service/vault"
	"log"
	"net"
	"os"
)

var (
	client *vault.ClientVault
)

const (

	// MongoURI environment variable
	MongoURI = "MONGO_URI"

	// DefaultMongoHost for running the app as a standalone server.
	DefaultMongoHost = "mongodb://localhost"

	// ENVIRONMENT OS variable
	ENVIRONMENT = "ENVIRONMENT"

	// DefaultEnvironment for running the app locally using docker-compose
	DefaultEnvironment = "default"
)

// MongoRepository type
type MongoRepository struct {
	*mgo.Session
}

// Repository interface
type Repository interface {
	Insert(database string, collection string, content interface{}) error
	FindAll(database string, collection string, query bson.M) ([]model.TagDAO, error)
	Find(database string, collection string, oid bson.ObjectId) (model.TagDAO, error)
	Delete(database string, collection string, oid bson.ObjectId) error
}

// NewRepository function to create an instance of Mongo repository
func NewRepository(config vault.Config) Repository {

	// default value
	uri := getDBURI(config)

	logger.Info.Printf("Initialising Mongo database session...")
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		log.Panicf("Failed to parse Mongo URI")
	}

	if requiresTLSConnection(uri) {
		handleTLS(dialInfo)
	}
	// get session
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		panic(err)
	}
	repository := &MongoRepository{session}
	return repository
}

// Implementation of Insert into Mongo repository
func (repo *MongoRepository) Insert(db string, collection string, body interface{}) error {
	return repo.Session.DB(db).C(collection).Insert(&body)
}

// Implementation of Find all from  Mongo repository for given id
func (repo *MongoRepository) FindAll(db string, collection string, query bson.M) ([]model.TagDAO, error) {
	var results []model.TagDAO
	err := repo.Session.DB(db).C(collection).Find(query).All(&results)
	return results, err
}

// Implementation of Insert into Mongo repository
func (repo *MongoRepository) Find(db string, collection string, oid bson.ObjectId) (model.TagDAO, error) {
	var result model.TagDAO
	err := repo.Session.DB(db).C(collection).FindId(oid).One(&result)
	return result, err
}

// Implementation of Delete
func (repo *MongoRepository) Delete(db string, collection string, oid bson.ObjectId) error {
	return repo.Session.DB(db).C(collection).RemoveId(oid)
}

// For `dev` and `prod` environment we enable TLS.
func requiresTLSConnection(uri string) bool {
	return os.Getenv(ENVIRONMENT) != DefaultEnvironment
}

// Add TLS configuration
func handleTLS(dialInfo *mgo.DialInfo) {
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
}

// Gets DB URI
func getDBURI(config vault.Config) string {
	// default uri
	uri := cfg.GetEnv(MongoURI, DefaultMongoHost)

	// override Mongo URI sourced from Vault if enabled.
	if config.Enabled {

		// Loads the secrets from vault server
		client, err := vault.NewClient(config)
		if err != nil {
			logger.Error.Printf("Failed to instantiate an instance of the vault client")
			panic(err)
		}
		uri = client.ReadSecret(MongoURI)
	}
	return uri
}
