package core

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/EricQAQ/Nebula/config"

	log "github.com/sirupsen/logrus"
)

func PrintRawInfo() {
	fmt.Println("Release Version: ", VERSION)
	fmt.Println("Golang Version: ", runtime.Version())
}

func PrintLogo() {
	log.Infof(`
████████╗██████╗  █████╗ ███████╗██████╗ 
╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔══██╗
   ██║   ██████╔╝███████║█████╗  ██║  ██║
   ██║   ██╔══██╗██╔══██║██╔══╝  ██║  ██║
   ██║   ██║  ██║██║  ██║███████╗██████╔╝
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═════╝ 
`)
}

func PrintInfo() {
	log.Infof("Welcome to Nebula.")
	log.Infof("Release Version: %s %s %s.", VERSION, runtime.GOOS, runtime.GOARCH)
	log.Infof("GoVersion: %s.", runtime.Version())
	configJSON, err := json.MarshalIndent(config.GetNebulaConfig(), "", "    ")
	if err != nil {
		panic(err)
	}
	log.Infof("******************Config******************")
	log.Infof("Config: \n%s", configJSON)
	log.Infof("******************************************")
}
