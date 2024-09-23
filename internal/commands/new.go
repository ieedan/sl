package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/game"
	"github.com/ieedan/sl/internal/util"
)

func (n *NewCmd) Run() error {
	db := database.Connect()
	defer db.Close()

	var existingName string

	err := db.Get(&existingName, "SELECT Name FROM Games WHERE Name = ?", n.Name)

	if err == nil && existingName != "" {
		log.Fatalf("You already have a game called %v", n.Name)
	}

	reader := bufio.NewReader(os.Stdin)

	// trainer and starter
	trainers := make(map[string]string)

	for {
		fmt.Println("Enter trainer name:")

		trainerName, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		trainerName = strings.TrimSpace(trainerName)

		if trainerName == "" {
			fmt.Printf("Trainer name cannot be blank!\n")
			continue
		}

		if _, ok := trainers[trainerName]; ok {
			fmt.Printf("You already added %v!\n", trainerName)
			continue
		}

		trainers[trainerName] = ""

		if len(trainers) >= 2 {
			fmt.Printf("Add more trainers? (y/N)")

			answer, err := reader.ReadString('\n')

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}

			answer = strings.ToLower(strings.TrimSpace(answer))

			if answer == "n" || answer == "" {
				break
			}
		}
	}

	for k := range trainers {
		for {
			fmt.Printf("Enter %v's starter:\n", k)

			pokemonName, err := reader.ReadString('\n')

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}

			pokemonName = strings.TrimSpace(pokemonName)

			if pokemonName == "" {
				fmt.Printf("Pokemon name cannot be blank!")
				continue
			}

			trainers[k] = pokemonName

			break
		}
	}

	fmt.Printf("Setting up new Soul Link %v..\n", n.Name)

	fmt.Println("Creating game...")

	result, err := db.Exec("INSERT INTO Games (Name) VALUES (?)", n.Name)
	if err != nil {
		log.Fatal(err)
	}

	gameId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	result, err = db.Exec(
		"INSERT INTO Routes (GameId, Name) VALUES (?, 'Starter')",
		gameId,
	)
	if err != nil {
		log.Fatal(err)
	}

	routeId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	for trainer, pokemon := range trainers {
		fmt.Printf("Setting up trainer %v...\n", trainer)

		result, err := db.Exec(
			"INSERT INTO Trainers (GameId, Name) VALUES (?, ?)",
			gameId,
			trainer,
		)
		if err != nil {
			log.Fatal(err)
		}

		trainerId, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(
			"INSERT INTO Pokemon (RouteId, TrainerId, Name) VALUES (?, ?, ?)",
			routeId,
			trainerId,
			pokemon,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	game.Play(n.Name)

	return nil
}

type NewCmd struct {
	Name string `arg:"" name:"name" help:"Name of the new game" type:"string"`
}
