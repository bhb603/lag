package cmd

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "lag",
		Short: "a server that lags",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

// Entrypoint for configuration, which will run before each command.
// Precedence:
//   0. explicit call to viper.Set
//   1. flag
//   2. env
//   3. config file (if provided)
func initConfig() {
	if len(cfgFile) > 0 {
		// if a cfgFile was provided via flag, load that
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Send()
		}
	} else {
		// otherwise try some default locations
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/lag")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("no config files found")
		}
	}
	if fileUsed := viper.ConfigFileUsed(); fileUsed != "" {
		log.Printf("config file: %s", viper.ConfigFileUsed())
	}
	viper.SetEnvPrefix("lag")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
}
