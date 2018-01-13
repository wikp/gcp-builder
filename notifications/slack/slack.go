package slack

import (
	"github.com/nlopes/slack"
	"github.com/wendigo/gcp-builder/context"
	"github.com/wendigo/gcp-builder/project"
	"log"
	"os"
)

var emptyParams = context.Params{}
var emptyAttachments []slackAttachment

const colorOK = "#00cd66"
const colorError = "#c42025"
const colorInfo = "#1e90ff"

type SlackNotificationProvider struct {
	channelId       string
	client          *slack.Client
	botName         string
	logger          *log.Logger
	params          context.Params
	threadTimestamp string
}

type slackAttachment struct {
	header  string
	content string
	color   string
}

func envOrDefault(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}

	return defaultValue
}

func NewSlackProvider(params context.Params) *SlackNotificationProvider {
	if token, exists := os.LookupEnv("SLACK_TOKEN"); exists {
		return &SlackNotificationProvider{
			client:    slack.New(token),
			channelId: envOrDefault("SLACK_CHANNEL_ID", "release"),
			botName:   envOrDefault("SLACK_BOT_NAME", "gcp-builder"),
			logger: log.New(
				os.Stdout, "[slack] ", log.Lmicroseconds,
			),
			params: params,
		}
	}

	return nil
}

func (s *SlackNotificationProvider) OnReleaseStarted() {
	s.send(
		":rocket: *{{ .ProjectFullName }}* is being released on *{{ .Environment }}* with version `{{ .ProjectVersion }}` :see_no_evil:",
		buildAttachment(s.params),
		emptyParams,
	)
}

func (s *SlackNotificationProvider) OnReleaseCompleted(err error) {

	if err == nil {
		s.send(
			`:rocket: Release has ended *successfully* :heart:`,
			emptyAttachments,
			emptyParams,
		)
	} else {
		s.send(
			`:rocket: Release has *failed* :cry:`,
			errorAttachment(err),
			emptyParams,
		)
	}
}

func (s *SlackNotificationProvider) OnImageBuilding(image project.Image) {

	merged := s.params.Merge(context.FromImage(image))

	s.send(
		merged.ExpandTemplate("Container *{{ .ImageName }}* is being built"),
		imageAttachment(merged),
		emptyParams,
	)
}

func (s *SlackNotificationProvider) OnImageBuilded(image project.Image, err error) {
	merged := s.params.Merge(context.FromImage(image))

	if err != nil {
		s.send(
			merged.ExpandTemplate("Container *{{ .ImageName }}* failed to build :cry:"),
			errorAttachment(err),
			emptyParams,
		)
	} else {
		s.send(
			merged.ExpandTemplate("Container *{{ .ImageName }}* was built successfully :grin:"),
			emptyAttachments,
			emptyParams,
		)
	}
}

func (s *SlackNotificationProvider) OnImagePushing(image project.Image) {
	merged := s.params.Merge(context.FromImage(image))

	s.send(
		merged.ExpandTemplate("Container {{ .ImageName }} is being pushed :boat:"),
		emptyAttachments,
		emptyParams,
	)
}

func (s *SlackNotificationProvider) OnImagePushed(image project.Image, err error) {
	merged := s.params.Merge(context.FromImage(image))

	if err != nil {
		s.send(
			merged.ExpandTemplate("Container *{{ .ImageName }}* failed to push to registry :cry:"),
			errorAttachment(err),
			emptyParams,
		)
	} else {
		s.send(
			merged.ExpandTemplate("Container *{{ .ImageName }}* was successfully pushed to registry :grin:"),
			emptyAttachments,
			emptyParams,
		)
	}
}

func (s *SlackNotificationProvider) OnConfigurationValidated(err error) {
	if err != nil {
		s.send(
			s.params.ExpandTemplate("Kubernetes deployment configuration is invalid :cry:"),
			errorAttachment(err),
			emptyParams,
		)
	} else {
		s.send(
			s.params.ExpandTemplate("Kubernetes deployment configuration is valid :small_airplane:"),
			emptyAttachments,
			emptyParams,
		)
	}
}

func (s *SlackNotificationProvider) OnDeploying() {
	s.send(
		s.params.ExpandTemplate("Deploying to *{{ .Environment }}* cluster *{{ .KubernetesCluster }}* :rocket:"),
		projectAttachment(s.params),
		emptyParams,
	)
}

func (s *SlackNotificationProvider) OnDeployed(err error) {
	if err != nil {
		s.send(
			s.params.ExpandTemplate("Failed to deploy to *{{ .EnvironmentName }}* :tired_face:"),
			errorAttachment(err),
			emptyParams,
		)
	} else {
		s.send(
			s.params.ExpandTemplate("Deployed successfully to *{{ .EnvironmentName }}* :trophy:"),
			emptyAttachments,
			emptyParams,
		)
	}
}

func (s *SlackNotificationProvider) IsConfigured() bool {
	return s.channelId != ""
}

func (s *SlackNotificationProvider) sendNotification(message string, attachments []slack.Attachment) error {
	parameters := slack.PostMessageParameters{
		Username:        s.botName,
		AsUser:          false,
		IconURL:         "http://lorempixel.com/48/48",
		ThreadTimestamp: s.threadTimestamp,
	}

	parameters.Attachments = attachments

	channelID, timestamp, err := s.client.PostMessage(s.channelId, message, parameters)

	if s.threadTimestamp == "" {
		s.threadTimestamp = timestamp
	}

	s.logger.Printf("Sent notification to channel: %s on %s with err: %v", channelID, timestamp, err)

	return err
}

func (s *SlackNotificationProvider) send(msg string, attachments []slackAttachment, params context.Params) {
	merged := params.Merge(s.params)

	slackAttachments := make([]slack.Attachment, 0)

	for _, attachment := range attachments {
		slackAttachments = append(slackAttachments, slack.Attachment{
			Pretext:    merged.ExpandTemplate(attachment.header),
			Text:       merged.ExpandTemplate(attachment.content),
			Color:      attachment.color,
			MarkdownIn: []string{"text", "pretext"},
		})
	}

	if err := s.sendNotification(merged.ExpandTemplate(msg), slackAttachments); err != nil {
		s.logger.Printf("Could not send slack notification: %v", msg)
	}
}
