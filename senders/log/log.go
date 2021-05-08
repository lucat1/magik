package log

import (
	"log"
)

type LogEmail struct {
  logger *log.Logger
}

func (s LogEmail) Send(to, title, body string) error {
  s.logger.Printf("Sending email to %s\nTitle:%s\nBody:%s", to, title, body)
  return nil
}

func NewSender(logger *log.Logger) (LogEmail, error) {
	return LogEmail{
    logger: logger,
	}, nil
}
