package notifications

import (
	"bytes"
	"github.com/nlopes/slack"
	"github.com/wendigo/gcp-builder/context"
	"html/template"
	"log"
	"os"
)

const defaultSlackTemplate = `:rocket: *[{{ .Environment }} {{ .Step }}]* {{ .ProjectFullName }}: {{ .Message }}`

const buildEnvironmentTemplate = `Platform: {{ .BuildPlatform }}
Build Url: {{ .BuildUrl }}
Repository: {{ .BuildRepository }}
Branch: {{ .BuildBranch }}
Tag: {{ .BuildTag }}
Resolved version: {{ .BuildVersion }}
`

type SlackNotificationProvider struct {
	channelId string
	client    *slack.Client
	botName   string
	logger    *log.Logger
	template  string
	params    context.Params
}

func envOrDefault(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}

	return defaultValue
}

func NewSlackProvider(params context.Params) NotificationsProvider {
	if token, exists := os.LookupEnv("SLACK_TOKEN"); exists {
		return &SlackNotificationProvider{
			client:    slack.New(token),
			channelId: envOrDefault("SLACK_CHANNEL_ID", "release"),
			botName:   envOrDefault("SLACK_BOT_NAME", "gcp-builder"),
			logger: log.New(
				os.Stdout, "[slack] ", log.Lmicroseconds,
			),
			template: envOrDefault("SLACK_MESSAGE_TEMPLATE", defaultSlackTemplate),
			params:   params,
		}
	}

	return nil
}

func (s *SlackNotificationProvider) SendNotification(text string, params context.Params) error {
	ctx := s.params.Merge(params).Merge(context.Params{
		"Message": text,
	})

	message, err := s.expand(s.template, ctx)
	if err != nil {
		return err
	}

	p := slack.PostMessageParameters{
		Username: s.botName,
		AsUser:   false,
		IconURL:  "http://lorempixel.com/48/48",
	}

	if buildInfo, err := s.expand(buildEnvironmentTemplate, ctx); err == nil {
		attachment := slack.Attachment{
			Pretext:    "",
			Text:       buildInfo,
			MarkdownIn: []string{"text", "pretext"},
			Color:      "#36a64f",
		}

		p.Attachments = []slack.Attachment{attachment}
	}

	channelID, timestamp, err := s.client.PostMessage(s.channelId, message, p)

	s.logger.Printf("Sent notification to channel: %s on %s with err: %v", channelID, timestamp, err)

	return err
}

func (s *SlackNotificationProvider) expand(tpl string, params context.Params) (string, error) {
	tmpl, err := template.New("slack-template").Parse(tpl)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}

	if err := tmpl.Execute(buffer, params); err != nil {
		return "", err
	}

	return string(buffer.String()), nil
}

func (s *SlackNotificationProvider) IsConfigured() bool {
	return s.channelId != ""
}
