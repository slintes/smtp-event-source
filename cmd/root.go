package cmd

import (
	"encoding/json"
	"github.com/slintes/smtpServer/internal/config"
	"github.com/slintes/smtpServer/internal/smtp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug        bool
	configString string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&configString, "config", "", "")
}

var rootCmd = &cobra.Command{
	Use:   "smtpServer",
	Short: "starts the smtp server",
	RunE: func(cmd *cobra.Command, args []string) error {

		//log.SetFormatter(&log.JSONFormatter{})
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Info("debug log on")
		}

		var cfg config.Config
		if err := json.Unmarshal([]byte(configString), &cfg); err != nil {
			log.Errorf("could not parse config: %s", configString)
			panic(err)
		}

		done := make(chan bool)
		log.Infof("debug log: %v", debug)

		err := smtp.Start(cfg)
		if err != nil {
			return err
		}

		// wait forever
		<-done
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("something went terrible wrong: %v", err)
	}
}
