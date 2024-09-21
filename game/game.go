package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"util"

	"database"
	tm "github.com/buger/goterm"
)

func Play(name string) {
	db := database.Connect()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		tm.Flush()
		tm.Clear()
		tm.MoveCursor(1, 1)

		game, exists := database.GetGame(db, name)

		if !exists {
			fmt.Printf("Couldn't find game with the name: '%v'", name)
			return
		}

		fmt.Println(game.String())

		if !game.IsDead() {
			fmt.Println("Waiting for command (catch, kill, quit)...")

			command, err := reader.ReadString('\n')

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}

			command = strings.ToLower(strings.TrimSpace(command))

			if command == "quit" {
				break
			}

			switch command {
			case "kill":
				kill(game)
			case "catch":
				catch(game)
			default:
				fmt.Println("Invalid command! Please enter a valid command (catch, kill, quit)")
			}
		} else {
			fmt.Println("Game over... Type `delete` to remove this game.")

			command, err := reader.ReadString('\n')

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}

			command = strings.ToLower(strings.TrimSpace(command))

			if command == "delete" {
				delete(game)
				break
			}
		}
	}
}

func delete(game *database.Game) {
	db := database.Connect()

	_, err := db.Exec("DELETE FROM Games WHERE Id = ?", game.Id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted '%v'!\n", game.Name)
}

func kill(game *database.Game) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter a route to kill:")

	route, err := reader.ReadString('\n')

	if util.IsCancel(err) {
		fmt.Println("Canceled.")
		os.Exit(0)
	}

	route = strings.TrimSpace(route)

	id, ok := game.GetRoute(route)

	if !ok {
		fmt.Printf("The route '%v' does not exist\n", route)
		return
	}

	db := database.Connect()

	_, err = db.Exec("UPDATE Routes SET PokemonAreAlive = 0 WHERE Id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func catch(game *database.Game) {
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

	for i, trainer := range *game.Trainers {
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
