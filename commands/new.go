package commands

import (
	"bufio"
	"database"
	"fmt"
	"game"
	"log"
	"os"
	"strings"
)

func (n *NewCmd) Run() error {
	db := database.Connect()
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM Games WHERE Name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(n.Name)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		log.Fatalf("You already have a game called %v", n.Name)
	}

	reader := bufio.NewReader(os.Stdin)

	trainers := make(map[string]string)

	for {
		fmt.Println("Enter trainer name:")

		trainerName, _ := reader.ReadString('\n')
	
		trainerName = strings.TrimSpace(trainerName)

		if _, ok := trainers[trainerName]; ok {
			fmt.Printf("You already added %v!\n", trainerName)
			continue
		}

		trainers[trainerName] = trainerName

		if (len(trainers) >= 2) {
			fmt.Printf("Add more trainers? (y/N)");

			answer, _ := reader.ReadString('\n');

			answer = strings.ToLower(strings.TrimSpace(answer))

			if answer == "n" || answer == "" {
				break
			}
		}
	}

	// need to ask for starters next

	fmt.Printf("Setting up new Soul Link %v..\n", n.Name)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = tx.Prepare("INSERT INTO Games (Name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.Name)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	game.Play(n.Name)

	return nil
}

type NewCmd struct {
	Name string `arg:"" name:"name" help:"Name of the new game" type:"string"`
}
