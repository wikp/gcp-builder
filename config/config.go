package config

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Steps         []string `arg:"positional,required"`
	KeyFile       string   `arg:"--key-file" help:"service key token file"`
	Branch        string   `arg:"--branch" help:"Current branch"`
	CommitSha     string   `arg:"--commit-sha" help:"Commit SHA"`
	Environment   string   `arg:"--env" help:"Current environment"`
	ProjectConfig string   `arg:"--config" help:"Project config yaml file"`
}

func Get() (*Args, error) {
	args := &Args{}

	args.Steps = []string{"all"}
	args.KeyFile = "key.json"
	args.ProjectConfig = "project.yml"
	args.Environment = "test"
	args.Branch = "master"
	args.CommitSha = ""

	arg.MustParse(args)

	return args, nil
}
