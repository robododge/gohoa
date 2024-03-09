package gohoa

import (
	"log"
	"sync"

	"github.com/JeremyLoy/config"
)

type MainConfig struct {
	MongoDBUrl           string `config:"MONGODB_URI"`
	MongoDBName          string `config:"MONGODB_NAME"`
	SlimMembersJson      string `config:"SLIM_MEMBERS_JSON"`
	SlimMembersJsonReval string `config:"SLIM_MEMBERS_JSON_REVAL"`

	StreetsJson string `config:"STREETS_JSON"`
}

var (
	mConfig    MainConfig
	configOnce sync.Once
)

func GetConfig() MainConfig {

	configOnce.Do(func() {
		err := config.From(".env.local").FromEnv().To(&mConfig)
		if err != nil {
			log.Fatal("Error loading config: ", err)
		}
	})
	return mConfig
}
