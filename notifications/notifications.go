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
	OnReleaseStarted()
	OnReleaseCompleted(error)
	OnImageBuilding(project.Image)
	OnImageBuilded(project.Image, error)
	OnImagePushing(project.Image)
	OnImagePushed(project.Image, error)
	OnConfigurationValidated(error)
	OnDeploying()
	OnDeployed(error)
	IsConfigured() bool
}

type DiscardingProvider struct {
}

func (d DiscardingProvider) OnReleaseStarted()                   {}
func (d DiscardingProvider) OnReleaseCompleted(error)            {}
func (d DiscardingProvider) OnImageBuilding(project.Image)       {}
func (d DiscardingProvider) OnImageBuilded(project.Image, error) {}
func (d DiscardingProvider) OnImagePushing(project.Image)        {}
func (d DiscardingProvider) OnImagePushed(project.Image, error)  {}
func (d DiscardingProvider) OnConfigurationValidated(error)      {}
func (d DiscardingProvider) OnDeploying()                        {}
func (d DiscardingProvider) OnDeployed(error)                    {}
func (d DiscardingProvider) IsConfigured() bool {
	return true
}
