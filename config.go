package gohoa

import (
	"fmt"
	"log"
	"os"
	"strings"
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

		var configBuilder *config.Builder

		log.Println(" *--* RAW ENV *--*")
		osEnvs := os.Environ()
		for _, env := range osEnvs {
			if strings.HasPrefix(env, "MONGODB_URI") {
				env = fmt.Sprintf("%s%s", env[:strings.Index(env, "=")+1], "*****")
			}
			log.Println(" *--* ENV: ", env)
		}

		currEnv, found := os.LookupEnv("ENV_MODE")
		if found && currEnv == "production" {
			log.Println(" *--* Running in PRODUCTION mode, trying to load any .env.production file.")
			configBuilder = config.FromOptional(".env.production")
		} else {
			log.Println(" *--* Running in DEVELOPMENT mode")
			configBuilder = config.FromOptional(".env").FromOptional(".env.local")
		}

		err := configBuilder.FromEnv().To(&mConfig)
		if err != nil {
			log.Fatal("Error loading config: ", err)
		}
	})
	return mConfig
}

// func getExistingFiles( filenames ...string) []string {
// 	var foundFileNames []string
// 	for _, filename := range filenames {
// 		if _, err := os.Stat(filename); err == nil {
// 			foundFileNames = append(foundFileNames, filename)
// 		}
// 	}
// 	return foundFileNames
// }
