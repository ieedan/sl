package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/util"
)

var Catch = Cmd{Name: "catch", Description: "Walks you through catching a new Pokemon", Run: catch}

func catch(args []string, game *database.Game) {
	reader := bufio.NewReader(os.Stdin)

	var routeName string

	for {
		fmt.Println("Enter the name of the route:")

		res, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		routeName = strings.TrimSpace(res)

		if routeName == "" {
			fmt.Println("You must enter the route name!")
			continue
		}

		break
	}

	pokemon := make(map[int64]string)

	for i, trainer := range game.Trainers {
		fmt.Printf("Enter %v's new pokemon:\n", trainer.Name)

		pokemonName, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		pokemonName = strings.TrimSpace(pokemonName)

		if pokemonName == "" {
			fmt.Println("You must enter a Pokemon!")
			i--
			continue
		}

		pokemon[trainer.Id] = pokemonName
	}

	db := database.Connect()
	defer db.Close()

	result, err := db.Exec(
		"INSERT INTO Routes (GameId, Name) VALUES (?, ?)",
		game.Id,
		routeName,
	)
	if err != nil {
		log.Fatal(err)
	}

	routeId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	for trainerId, pokemonName := range pokemon {
		_, err = db.Exec(
			"INSERT INTO Pokemon (RouteId, TrainerId, Name) VALUES (?, ?, ?)",
			routeId,
			trainerId,
			pokemonName,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
