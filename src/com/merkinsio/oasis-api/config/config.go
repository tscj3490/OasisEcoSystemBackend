package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config Configuration instance from Viper
var Config *viper.Viper

func init() {
	Config = viper.New()
	configPath := os.Getenv("CONFIG_FILE")

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Need to defined CONFIG_FILE environment variable")
	}

	Config.SetConfigFile(configPath)

	err := Config.ReadInConfig()
	if err != nil {
		// log.LogFile.WithError(err).Fatalf("Error reading the configuration file `%s`", configPath)
		panic(err.Error())
	}

	Config.SetEnvPrefix("APP")
	replacer := strings.NewReplacer(".", "_")
	Config.SetEnvKeyReplacer(replacer)
	Config.AutomaticEnv()

	fmt.Printf("Using config: %s\n", Config.ConfigFileUsed())
}
