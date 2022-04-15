package main

import (
	"embed"

	"github.com/wewillapp-com/we-cli/cmd"
)

//go:embed templates
var t embed.FS

//go:embed version.txt
var v string

func main() {
	//assign embed templates to cmd
	cmd.TemplateFS = t
	cmd.CurrentVersion = v
	cmd.Execute()
}
