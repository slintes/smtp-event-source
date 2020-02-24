package smtp

import (
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/flashmob/go-guerrilla/response"
	"github.com/slintes/smtpServer/internal/config"
)

type ReceiverValidatorProcessor struct {
	cfg config.Config
}

func NewReceiverValidatorProcessor(cfg config.Config) *ReceiverValidatorProcessor {
	return &ReceiverValidatorProcessor{
		cfg: cfg,
	}
}

func (r *ReceiverValidatorProcessor) Decorater() backends.Decorator {
	return func(p backends.Processor) backends.Processor {
		return backends.ProcessWith(
			func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {

				if task == backends.TaskValidateRcpt {

					for _, rec := range e.RcptTo {
						if m := config.FindMappingForRecipient(rec.User, r.cfg.Mappings); m != nil {
							return p.Process(e, task)
						}
					}

					backends.Log().WithError(backends.NoSuchUser).Info("invalid user")
					return backends.NewResult(response.Canned.FailRcptCmd), backends.NoSuchUser

				} else if task == backends.TaskSaveMail {

					return p.Process(e, task)

				}

				return p.Process(e, task)

			},
		)
	}
}
