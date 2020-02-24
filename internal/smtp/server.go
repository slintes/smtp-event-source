package smtp

import (
	"fmt"
	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/slintes/smtpServer/internal/config"
	"github.com/slintes/smtpServer/internal/event"
)

func Start(cfg config.Config) error {

	smtpCfg := &guerrilla.AppConfig{
		Servers: []guerrilla.ServerConfig{
			{
				ListenInterface: "0.0.0.0:2525",
				XClientOn:       true,
				IsEnabled:       true,
			},
		},
		AllowedHosts: []string{cfg.HostName},
		LogFile:      "stdout",
		LogLevel:     "debug",
		BackendConfig: backends.BackendConfig{
			"save_workers_size": 1,
			"save_process":      "HeadersParser|Debugger|Hasher|Header|Webhook",
			"validate_process":  "ReceiverValidator",
		},
	}

	d := guerrilla.Daemon{
		Config: smtpCfg,
	}

	d.AddProcessor("ReceiverValidator", NewReceiverValidatorProcessor(cfg).Decorater)
	d.AddProcessor("Webhook", event.NewWebhookProcessor(cfg).Decorater)

	err := d.Start()
	if err != nil {
		return err
	}

	fmt.Println("SMTP server Started!")
	return nil
}
