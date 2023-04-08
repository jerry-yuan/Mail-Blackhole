package storage

import (
	"fmt"

	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/sirupsen/logrus"
)

// Storage represents a storage backend
type Storage interface {
	Store(m *data.Message) (string, error)
	List(start, limit int) (*data.Messages, error)
	Search(kind, query string, start, limit int) (*data.Messages, int, error)
	Count() int
	DeleteOne(id string) error
	DeleteAll() error
	Load(id string) (*data.Message, error)
}

func Create(conf config.StorageConfig) (msgStorage Storage, err error) {
	switch conf.StorageType {
	case "memory":
		logrus.Infof("Using in-memory storage")
		msgStorage = CreateInMemory()
	case "mongodb":
		logrus.Infof("Using MongoDB message storage")
		msgStorage, err = CreateMongoDB(conf.MongoURI, conf.MongoDb, conf.MongoColl)
	case "maildir":
		logrus.Println("Using maildir message storage")
		msgStorage = CreateMaildir(conf.MaildirPath)
	default:
		return nil, fmt.Errorf("invalid storage type %s", conf.StorageType)
	}
	return
}
