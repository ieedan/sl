package main

import (
	"commands"
	"github.com/alecthomas/kong"
)

var CLI struct {
	New commands.NewCmd `cmd:"" help:"Create new soul link."`

	List commands.ListCmd `cmd:"" help:"List existing games."`

	Resume commands.ResumeCmd `cmd:"" help:"Pick up where you left off on an existing soul link."`
}

func main() {
	ctx := kong.Parse(&CLI)

	err := ctx.Run()

	ctx.FatalIfErrorf(err)
}
