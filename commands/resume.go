package commands

import (
	"fmt"
)

func (n *ResumeCmd) Run() error {
	fmt.Println("resume", n.Name)

	return nil
}

type ResumeCmd struct {
	Name string `arg:"" optional:"" name:"name" help:"Name of the game to resume"`
}