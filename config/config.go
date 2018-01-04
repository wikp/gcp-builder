package config

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Steps         []string `arg:"positional,required"`
	KeyFile       string   `arg:"--key-file" help:"service key token file"`
	Environment   string   `arg:"--env" help:"Current environment"`
	ProjectConfig string   `arg:"--config" help:"Project config yaml file"`
}

func Get() (*Args, error) {
	args := &Args{}

	args.Steps = []string{"all"}
	args.KeyFile = "key.json"
	args.ProjectConfig = "project.yml"
	args.Environment = "test"

	arg.MustParse(args)

	return args, nil
}
