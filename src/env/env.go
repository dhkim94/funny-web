package env

import (
	"github.com/spf13/viper"
	"fmt"
	"log"
	"cklog"
)

var clog *cklog.Cklogger
var conf *viper.Viper

func Init(confDir string, confFile string) {

	conf = viper.New()
	conf.SetConfigName(confFile)
	conf.AddConfigPath(confDir)

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalf("can't read config file [%s/%s.properties]",
			confDir, confFile)
	}
	fmt.Printf("read config file [%s/%s.properties]\n", confDir, confFile);

	logLevel := fmt.Sprintf("%s", conf.Get("log.level"))
	logOut := fmt.Sprintf("%s", conf.Get("log.output"))

	if logOut == "file" && conf.Get("log.file") == nil {
		logOut = "stdout"
	}

	clog = cklog.NewLogger(logLevel, logOut,
		fmt.Sprintf("%s", conf.Get("log.file")))
}

func GetLogger() *cklog.Cklogger {
	return clog
}

func GetConf() *viper.Viper {
	return conf
}

func GetConfig(key string) string {
	return fmt.Sprintf("%s", conf.Get(key))
}
