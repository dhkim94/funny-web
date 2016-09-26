package main

import (
	"fmt"
	ev "env"
)


const (
	configFile	= "config"
	configDir	= "/Users/dhkim/Documents/Develop/project.my/funny-web"
)

func main() {
	ev.Init(configDir, configFile)


	fmt.Println("----- main")

	//tm := time.Now()
	//ctm := cktime.NewCktime(tm, "YYYYMMDD")
	//fmt.Println(ctm.ToString())

	ev.LOG.Info("----info log in main 123 [%d]", 1)

}
