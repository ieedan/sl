package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/util"
	"github.com/jmoiron/sqlx"
)

var Kill = Cmd{Name: "kill", Description: "Kill all Pokemon in a route", Args: []Arg{
	{Name: "route", Description: "Name of the route to kill (case insensitive)", Optional: true},
}, Run: kill}

func kill(args []string, game *database.Game) {
	var route string

	if nameOfRoute != nil && *nameOfRoute != "" {
		route = *nameOfRoute
	}

	if route == "" {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Please enter a route to kill:")

		r, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		trimmed := strings.TrimSpace(r)

		route = trimmed
	}

	id, ok := game.GetRoute(route)

	if !ok {
		fmt.Printf("The route '%v' does not exist\n", route)
		return
	}

	killRoutes(id)
}

func killRoutes(routeIds ...int64) {
	db := database.Connect()
	defer db.Close()

	query, args, err := sqlx.In("UPDATE Routes SET PokemonAreAlive = 0 WHERE Id IN (?)", routeIds)
	if err != nil {
		log.Fatal(err)
	}

	query = db.Rebind(query)

	_, err = db.Exec(query, args...)
	if err != nil {
		log.Fatal(err)
	}
}
