package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config holds user settings, some of which might be overridable by command-line flags.
type Config struct {
	Context   string
	Namespace string
}

// ParseConfig reads in user configuration from files, with some settings optionally being overridable via command-line flags.
func ParseConfig(cmd *cobra.Command) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configPath())
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.SetEnvPrefix("frink")
	v.AutomaticEnv()

	v.BindPFlags(cmd.Flags())

	if err := v.ReadInConfig(); err != nil {
		// TODO: Log ConfigFileNotFoundError when/if we implement logging?
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

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
