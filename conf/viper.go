package conf

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadViper(localFile string) *viper.Viper {
	confFile := strings.TrimSpace(os.Getenv("APP_CONF"))
	if len(confFile) == 0 {
		confFile = localFile
	}
	if len(confFile) == 0 {
		panic("conf is nil")
	}
	log.Println("load conf file:", confFile)

	conf := viper.New()
	conf.SetConfigFile(confFile)
	conf.SetConfigType("yaml")
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := conf.ReadInConfig()
	if err != nil {
		log.Panicf("read config err: %+v", err)
	}

	return conf
}
