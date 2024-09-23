package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/table"
)

func (l *ListCmd) Run() error {
	db := database.Connect()
	defer db.Close()

	var names []string

	err := db.Select(&names, "Select Name FROM Games G")
	if err != nil {
		log.Fatal(err)
	}

	games := []database.Game{}

	for _, name := range names {
		// we know that ok will always be true since we just got it
		game, _ := database.GetGame(db, name)

		games = append(games, *game)
	}

	t := table.New(table.DEFAULT_OPTIONS)

	t.AddHeader("Name", "Trainers", "Last Route", "Remaining Pokemon", "Started")

	for _, game := range games {
		lastRoute := game.Routes[len(game.Routes)-1]

		remaining := 0

		for _, route := range game.Routes {
			if route.PokemonAreAlive {
				remaining += 1
			}
		}

		columns := []string{
			game.Name,
			strconv.Itoa(len(game.Trainers)),
			lastRoute.Name,
			fmt.Sprintf("%v", remaining),
			game.CreatedAt.Local().Format("01-02-2006 03:04 PM"),
		}

		t.AddRow(columns...)
	}

	fmt.Println(t.String())

	return nil
}

type ListCmd struct {
}
