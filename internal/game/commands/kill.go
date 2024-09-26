package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/util"
)

var Kill = Cmd{Name: "kill", Description: "Kill all Pokemon in a route", Args: []Arg{
	{Name: "route", Description: "Name of the route to kill (case insensitive)", Optional: true},
}, Run: kill}

func kill(args []string, game *database.Game) {
	var route string

	if len(args) > 1 {
		route = strings.Join(args[1:], " ")
	}

	for route == "" {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Please enter a route to kill:")

		r, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		trimmed := strings.TrimSpace(r)

		if route == "" {

		}

		route = trimmed
	}

	id, ok := game.GetRoute(route)

	if !ok {
		fmt.Printf("The route '%v' does not exist\n", route)
		return
	}

	KillRoutes(id)
}
