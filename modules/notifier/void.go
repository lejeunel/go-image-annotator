package notifier

import (
	"log/slog"
)

type VoidNotifier struct {
	slog.Logger
}

func (n VoidNotifier) Notify(recipient, url string) {
	n.Logger.Info("notifying password reset token", "recipient", recipient)
}

func (n VoidNotifier) Check() error {
	return nil
}
