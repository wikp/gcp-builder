package config

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Steps         []string `arg:"positional,required"`
	Environment   string   `arg:"--env" help:"Current environment"`
	ProjectConfig string   `arg:"--config" help:"Project config yaml file"`
	Update        bool     `arg:"--update" help:"Update gcloud components"`
}

func Get() (*Args, error) {
	args := &Args{}

	args.Steps = []string{"all"}
	args.ProjectConfig = "project.yml"
	args.Environment = "test"
	args.Update = false

	arg.MustParse(args)

	return args, nil
}
