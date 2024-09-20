package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"database"
)

func Play(name string) {
	db := database.Connect()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	game := database.GetGame(db, name)

	fmt.Printf("%v", game)

	for {
		fmt.Println("Waiting for command (catch, kill, quit)...")

		command, _ := reader.ReadString('\n')

		command = strings.TrimSpace(command)

		if command == "quit" {
			break
		}

		switch command {
		case "kill":
			kill()
		case "catch":
			catch()
		default:
			fmt.Println("Invalid command! Please enter a valid command (catch, kill, quit)")
		}


	}
}

func kill() {

}

func catch() {

}