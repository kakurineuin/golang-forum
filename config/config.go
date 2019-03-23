package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Viper 提供取得設定檔中的值的方法。
var Viper viper.Viper

// Init 初始化設定檔。
func Init(configPath, configName string) {
	Viper = *viper.New()
	Viper.SetConfigName(configName)
	Viper.AddConfigPath(configPath)
	err := Viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// InitByEnv 從環境變數初始化設定檔。
func InitByEnv() {
	Viper = *viper.New()
	Viper.AutomaticEnv()
	Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
