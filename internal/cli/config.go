package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func configPath() string {
	configPath := os.Getenv("XDG_CONFIG_HOME")
	if configPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		configPath = path.Join(home, ".config")
	}

	frinkPath := path.Join(configPath, "frink")
	return frinkPath
}

// InitConfig wires up the default configuration backend.
func InitConfig() {
	viper.AddConfigPath(configPath())
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// TODO: Log ConfigFileNotFoundError when/if we implement logging?
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
