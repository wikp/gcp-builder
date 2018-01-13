package context

import (
	"github.com/wendigo/gcp-builder/kubernetes"
	"github.com/wendigo/gcp-builder/platforms"
)

type Params map[string]string

func (p Params) Merge(l Params) Params {
	r := Params{}

	for k, v := range p {
		r[k] = v
	}

	for k, v := range l {
		r[k] = v
	}

	return r
}

func From(context *kubernetes.Context, platform platforms.Platform) Params {
	return Params{
		"Environment":       context.CurrentEnvironment.Name,
		"KubernetesCluster": context.CurrentEnvironment.Kubernetes.Cluster,
		"KubernetesZone":    context.CurrentEnvironment.Kubernetes.Zone,
		"CloudProject":      context.CurrentEnvironment.Cloud.Project,
		"CloudRegistry":     context.CurrentEnvironment.Cloud.Registry,
		"BuildUrl":          platform.BuildUrl(),
		"BuildRepository":   platform.RepositoryUrl(),
		"BuildPlatform":     platform.Name(),
		"BuildNumber":       platform.CurrentBuildNumber(),
		"BuildBranch":       platform.CurrentBranch(),
		"BuildTag":          platform.CurrentTag(),
		"BuildVersion":      context.Version,
		"ProjectName":       context.Config.Project.Name,
		"ProjectDomain":     context.Config.Project.Domain,
		"ProjectContext":    context.Config.Project.Context,
		"ProjectVersion":    context.Version,
		"ProjectFullName":   context.Config.Project.FullName(),
	}
}
