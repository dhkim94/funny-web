package env

import (
	"github.com/spf13/viper"
	"fmt"
	"log"
	"cklog"
)

var LOG *cklog.Cklogger
var Viper *viper.Viper

func Init(confDir string, confFile string) {

	Viper = viper.New()
	Viper.SetConfigName(confFile)
	Viper.AddConfigPath(confDir)

	err := Viper.ReadInConfig()
	if err != nil {
		log.Fatalf("can't read config file [%s/%s.properties]",
			confDir, confFile)
	}
	fmt.Printf("read config file [%s/%s.properties]\n", confDir, confFile);

	logLevel := fmt.Sprintf("%s", Viper.Get("log.level"))
	logOut := fmt.Sprintf("%s", Viper.Get("log.output"))

	if logOut == "file" && Viper.Get("log.file") == nil {
		logOut = "stdout"
	}

	LOG = cklog.NewLogger(logLevel, logOut,
		fmt.Sprintf("%s", Viper.Get("log.file")))




}
