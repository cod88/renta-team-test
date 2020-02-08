package appconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	AppRoot string
	Appname string
}

var config AppConfig

func LoadConfig() {
	execFile, _ := os.Executable()
	approot := filepath.Dir(filepath.Dir(execFile))
	config.AppRoot = approot
	if _, err := toml.DecodeFile(approot+"/config/config.toml", &config); err != nil {
		fmt.Println("We have an error. ", err)
	}
}

func GetConfig() *AppConfig {
	if config == (AppConfig{}) {
		LoadConfig()
	}
	return &config
}
