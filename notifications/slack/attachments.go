package slack

import (
	"github.com/wendigo/gcp-builder/context"
	"fmt"
)

func errorAttachment(err error) []slackAttachment {
	return []slackAttachment{{
		header:  "Error details",
		content: err.Error(),
		color:   colorError,
	}}
}

func buildAttachment(ctx context.Params) []slackAttachment {
	template := `Build platform: *{{ .BuildPlatform }}*
Build url: {{ .BuildUrl }}
Repository: {{ .BuildRepository }}
Commit: *{{ .BuildCommit }}*
Branch: *{{if .BuildBranch }}{{ .BuildBranch }}{{else}}n/a{{end}}*
Tag: *{{if .BuildTag }}{{ .BuildTag }}{{else}}n/a{{end}}*
Version: *{{ .ProjectVersion }}*
Environment: *{{ .Environment }}*`

	return []slackAttachment{{
		header:  "",
		content: ctx.ExpandTemplate(template),
		color:   colorInfo,
	}}
}

func outputAttachment(output string) []slackAttachment {
	return []slackAttachment{{
		header:  "",
		content: fmt.Sprintf("```%s```", output),
		color:   colorOK,
	}}
}

func errorOutputAttachment(output string, err error) []slackAttachment {
	return append(outputAttachment(output), errorAttachment(err)...)
}

func projectAttachment(ctx context.Params) []slackAttachment {
	template := `Version: *{{ .ProjectVersion }}*
Environment: *{{ .Environment }}*
Cluster: *{{ .KubernetesCluster }}*
Zone: *{{ .KubernetesZone }}*
Project: *{{ .CloudProject }}*
Registry: *{{ .CloudRegistry }}*
`

	return []slackAttachment{{
		header:  "",
		content: ctx.ExpandTemplate(template),
		color:   colorInfo,
	}}
}

func imageAttachment(ctx context.Params) []slackAttachment {
	template := `Container name: *{{ .ImageName }}*
Registry: {{ .CloudRegistry }}/{{ .ProjectFullName}}/{{ .ImageName }}:{{ .BuildVersion }}`

	return []slackAttachment{{
		header:  "",
		content: ctx.ExpandTemplate(template),
		color:   colorInfo,
	}}
}
