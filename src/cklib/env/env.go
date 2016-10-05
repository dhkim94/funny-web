package env

import (
	"github.com/spf13/viper"
	"fmt"
	"cklib/cklog"
	"path/filepath"
	"strings"
)

var (
	clog *cklog.Cklogger
	conf *viper.Viper
)

func Init(confFileName string) bool {
	confDir := filepath.Dir(confFileName)
	baseName := strings.Split(filepath.Base(confFileName), ".")
	baseNameLen := len(baseName)

	if baseNameLen < 2 {
		fmt.Printf("[FAIL] Invalid properties file: %s\n", confFileName)
		return false
	}

	if baseName[baseNameLen - 1] != "properties" {
		fmt.Printf("[FAIL] config file is not properties file: %s\n", confFileName)
		return false
	}

	confFile := strings.Join(baseName[:baseNameLen - 1], ".")

	conf = viper.New()
	conf.SetConfigName(confFile)
	conf.AddConfigPath(confDir)

	err := conf.ReadInConfig()
	if err != nil {
		fmt.Printf("[FAIL] Not found config file: %s\n", confFileName)
		return false
	}
	//fmt.Printf("read config file [%s/%s.properties]\n", confDir, confFile);

	logLevel := fmt.Sprintf("%s", conf.Get("log.level"))
	logOut := fmt.Sprintf("%s", conf.Get("log.output"))

	if logOut == "file" && conf.Get("log.file") == nil {
		logOut = "stdout"
	}

	clog = cklog.NewLogger(logLevel, logOut,
		fmt.Sprintf("%s", conf.Get("log.file")))

	return true
}

func GetLogger() *cklog.Cklogger {
	return clog
}

func GetConf() *viper.Viper {
	return conf
}

func GetValue(key string) string {
	return fmt.Sprintf("%s", conf.Get(key))
}
