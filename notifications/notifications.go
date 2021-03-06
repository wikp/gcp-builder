package notifications

import (
	"github.com/wendigo/gcp-builder/context"
	"github.com/wendigo/gcp-builder/notifications/slack"
	"github.com/wendigo/gcp-builder/project"
)

func Get(params context.Params) NotificationsProvider {
	for _, provider := range []NotificationsProvider{
		slack.NewSlackProvider(params),
		DiscardingProvider{},
	} {
		if provider != nil && provider.IsConfigured() {
			return provider
		}
	}

	return nil
}

type NotificationsProvider interface {
	OnReleaseStarted([]string)
	OnReleaseCompleted([]string, error)
	OnImageBuilding(project.Image)
	OnImageBuilt(project.Image, string, error)
	OnImagePushing(project.Image)
	OnImagePushed(project.Image, string, error)
	OnConfigurationValidated(error)
	OnDeploying()
	OnDeployed(string, error)
	IsConfigured() bool
}

type DiscardingProvider struct {
}

func (d DiscardingProvider) OnReleaseStarted([]string)                  {}
func (d DiscardingProvider) OnReleaseCompleted([]string, error)         {}
func (d DiscardingProvider) OnImageBuilding(project.Image)              {}
func (d DiscardingProvider) OnImageBuilt(project.Image, string, error)  {}
func (d DiscardingProvider) OnImagePushing(project.Image)               {}
func (d DiscardingProvider) OnImagePushed(project.Image, string, error) {}
func (d DiscardingProvider) OnConfigurationValidated(error)             {}
func (d DiscardingProvider) OnDeploying()                               {}
func (d DiscardingProvider) OnDeployed(string, error)                   {}
func (d DiscardingProvider) IsConfigured() bool {
	return true
}
