package cfg

import (
	_ "embed"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/ul-gaul/go-basestation/utils"
)

const (
	Name      = "gaul-config"
	EnvPrefix = "GAUL"
)

//go:embed default-config.yaml
var defaultCfgBytes []byte

func Initialize(cfgFile string) {
	var err error

	// Load default config
	defaultCfgMap := make(map[string]interface{})
	utils.CheckErr(yaml.Unmarshal(defaultCfgBytes, defaultCfgMap))
	for key, value := range defaultCfgMap {
		viper.SetDefault(key, value)
	}

	// Load and merge config file in home, cwd and executable directories
	{
		var home, exe string

		// Find home directory.
		home, err = homedir.Dir()
		utils.CheckErr(err)

		// Find executable path
		exe, err = utils.GetExecutablePath()
		utils.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath(exe)

		viper.SetConfigName(Name)

		if err = viper.MergeInConfig(); err != nil {
			// Ignore not found error
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				utils.CheckErr(err)
			}
		}
	}

	// Load and merge specified config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		utils.CheckErr(viper.MergeInConfig())
	}

	// Load and merge config from environnement variables
	viper.SetEnvPrefix(EnvPrefix)
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	utils.CheckErr(viper.UnmarshalKey("Comms", &Comms))
	utils.CheckErr(viper.UnmarshalKey("Frontend", &Frontend))
}


// TODO validate settings