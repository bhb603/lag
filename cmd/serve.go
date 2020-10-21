package cmd

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bhb603/lag/server"
)

var serveCommand = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"server", "s"},
	Short:   "Launch the server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &server.Config{
			Port:        viper.GetString("server.port"),
			MaxDataSize: viper.GetString("server.max_data_size"),
		}

		if str := viper.GetString("server.max_lag"); str != "" {
			d, err := time.ParseDuration(str)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			cfg.MaxLag = d
		}

		server.Serve(cfg)
	},
}

func init() {
	serveCommand.Flags().String("port", "8080", "server port")
	viper.BindPFlag("server.port", serveCommand.Flags().Lookup("port"))
	rootCmd.AddCommand(serveCommand)
}
