package event

import (
	"encoding/base64"
	"github.com/slintes/smtpServer/internal/config"
	"golang.org/x/net/html/charset"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/flashmob/go-guerrilla/response"
	"github.com/jhillyerd/enmime"
)

type WebhookProcessor struct {
	cfg config.Config
}

func NewWebhookProcessor(cfg config.Config) *WebhookProcessor {
	return &WebhookProcessor{
		cfg: cfg,
	}
}

func (w *WebhookProcessor) Decorater() backends.Decorator {
	return func(p backends.Processor) backends.Processor {
		return backends.ProcessWith(
			func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {

				if task == backends.TaskSaveMail {

					var mapping *config.Mapping
					for _, rec := range e.RcptTo {
						if mapping = config.FindMappingForRecipient(rec.User, w.cfg.Mappings); mapping != nil {
							break
						}
					}

					if mapping == nil {
						backends.Log().WithError(backends.NoSuchUser).Info("invalid user, NOT CATCHED BY VALIDATION!?")
						return backends.NewResult(response.Canned.FailRcptCmd), backends.NoSuchUser
					}

					mimeEnvelope, err := enmime.ReadEnvelope(e.NewReader())
					if err != nil {
						backends.Log().WithError(backends.StorageError).Infof("can't parse message: %v", err)
						return backends.NewResult(response.Canned.FailBackendTransaction), backends.StorageError
					}

					for _, attachment := range mimeEnvelope.Attachments {

						re := regexp.MustCompile(mapping.FileNameRegEx)
						if !re.MatchString(strings.ToLower(attachment.FileName)) {
							backends.Log().Infof("skipping attachment: %v", attachment.FileName)
							continue
						}

						bytes := attachment.Content
						if mapping.FileEncoding != "" {
							enc, encName := charset.Lookup(mapping.FileEncoding)
							if enc != nil {
								dec := enc.NewDecoder()
								bytes, err = dec.Bytes(attachment.Content)
								if err == nil {
									backends.Log().Infof("decoded message using : %v", encName)
								}
							}
						}

						data := base64.StdEncoding.EncodeToString(bytes)
						res, err := http.Post(mapping.WebhookUrl, "text/plain", strings.NewReader(data))
						if err != nil || res.StatusCode >= 300 {
							status := ""
							if err != nil {
								status = err.Error()
							} else if res != nil && res.Status != "" {
								status = res.Status
							} else if res != nil {
								status = strconv.Itoa(res.StatusCode)
							}
							backends.Log().WithError(backends.StorageError).Infof("couldn't call webhook: %v", status)
							return backends.NewResult(response.Canned.FailBackendTransaction), backends.StorageError
						}

						// we only handle the 1st csv attachment for now
						break
					}

					return p.Process(e, task)

				} else if task == backends.TaskValidateRcpt {

					return p.Process(e, task)

				}

				return p.Process(e, task)

			},
		)
	}
}
