package game

import (
	"fmt"
	"util"
)

type Cmd struct {
	Name        string
	Description string
	Args        *[]Arg
}

type Arg struct {
	Name        string
	Description string
	Optional    bool
}

// Return help for a command
func (cmd *Cmd) Help() string {
	help := util.LPad(fmt.Sprintf("Usage: %v ", cmd.Name), 2)

	minArgLength := 0

	if cmd.Args != nil {
		for _, arg := range *cmd.Args {
			if len(arg.Name) > minArgLength {
				minArgLength = len(arg.Name)
			}

			if arg.Optional {
				help += fmt.Sprintf("<%v> ", arg.Name)
			} else {
				help += fmt.Sprintf("[%v] ", arg.Name)
			}
		}
	}

	help += "\n\n" + util.LPad(fmt.Sprintln(cmd.Description), 2) + "\n"

	if cmd.Args != nil {
		for _, arg := range *cmd.Args {
			if arg.Optional {
				help += util.LPad(fmt.Sprintf("%v%v\n", util.PadRightMin(arg.Name, minArgLength+4), arg.Description), 2)
			} else {
				help += util.LPad(fmt.Sprintf("%v%v\n", util.PadRightMin(arg.Name, minArgLength+4), arg.Description), 2)
			}
		}
	}

	return help
}

func Help(commands *[]Cmd) string {
	var help string

	strNames := util.Map(commands, func(cmd Cmd, i int) string {
		return cmd.Name
	})

	minNameLength := util.MinLength(&strNames)

	for _, command := range *commands {
		help += util.LPad(fmt.Sprintf("%v%v\n", util.PadRightMin(command.Name, minNameLength+4), command.Description), 2)
	}

	return help
}

var Commands = []Cmd{
	{Name: "catch", Description: "Walks you through catching a new Pokemon"},
	{Name: "kill", Description: "Kill all Pokemon in a route", Args: &[]Arg{
		{Name: "route", Description: "Name of the route to kill (case insensitive)", Optional: true},
	}},
	{Name: "end", Description: "Ends the game (whiteout / blackout)"},
	{Name: "quit", Description: "Quit from the terminal (ctrl + c)"},
	{Name: "help", Description: "Displays help", Args: &[]Arg{
		{Name: "command", Description: "Name of a command to display help for", Optional: true},
	}},
}
