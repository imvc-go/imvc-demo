package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
)

const (
	configKeyServerDev = "app.server.dev"
)

var instance *viper.Viper

func InitConfigWithPaths(paths ...string) {
	instance = viper.New()
	instance.SetConfigName("config")

	for _, path := range paths {
		instance.AddConfigPath(path)
	}

	instance.SetConfigType("yaml")

	// Find and read the config file
	if err := instance.ReadInConfig(); err != nil {
		panic(err)
	}

	instance.OnConfigChange(func(e fsnotify.Event) {
		golog.Info("Config reloading", e.Name)
	})

	instance.WatchConfig()
	//instance.Debug()

	golog.Info("Init global config")
}

func InitConfig() {
	InitConfigWithPaths("./configs/", "../configs/", "../../configs/", "../../../configs/")
}

func GetConfig() *viper.Viper {
	if instance == nil {
		InitConfig()
	}
	return instance
}

// IsDevelopment
func IsDebug() bool {
	return GetConfig().GetBool(configKeyServerDev)
}
