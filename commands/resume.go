package commands

import (
	"game"
)

func (n *ResumeCmd) Run() error {
	game.Play(n.Name)

	return nil
}

type ResumeCmd struct {
	Name string `arg:"" optional:"" name:"name" help:"Name of the game to resume"`
}