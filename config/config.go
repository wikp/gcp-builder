package config

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Steps         []string `arg:"positional,required"`
	Environment   string   `arg:"--env" help:"Current environment"`
	ProjectConfig string   `arg:"--config" help:"Project config yaml file"`
}

func Get() (*Args, error) {
	args := &Args{}

	args.Steps = []string{"all"}
	args.ProjectConfig = "project.yml"
	args.Environment = "test"

	arg.MustParse(args)

	return args, nil
}
