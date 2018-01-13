package notifications

import (
	"github.com/wendigo/gcp-builder/context"
)

func Get(params context.Params) NotificationsProvider {
	for _, provider := range []NotificationsProvider{
		NewSlackProvider(params),
	} {
		if provider != nil && provider.IsConfigured() {
			return provider
		}
	}

	return nil
}

type NotificationsProvider interface {
	SendNotification(text string, params context.Params) error
	IsConfigured() bool
}
