package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/ieedan/sl/internal/args"
	"github.com/ieedan/sl/internal/database"
	gc "github.com/ieedan/sl/internal/game/commands"
	"github.com/ieedan/sl/internal/util"

	tm "github.com/buger/goterm"
)

func Play(name string) {
	db := database.Connect()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	// we clear initially for when we come from a new game
	tm.Flush()
	tm.Clear()
	tm.MoveCursor(1, 1)

	for {
		// clear screen
		tm.Flush()
		tm.Clear()
		tm.MoveCursor(1, 1)

		// get game
		game, exists := database.GetGame(db, name)

		if !exists {
			fmt.Printf("Couldn't find game with the name: '%v'", name)
			return
		}

		// display game
		fmt.Println(game.String())

		// run end loop when game is dead
		if game.IsDead() {
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

			continue
		}

		commandList := strings.Join(util.Map(&gc.Commands, func(cmd gc.Cmd, i int) string {
			return cmd.Name
		}), ", ")

		fmt.Printf("Waiting for command (%v)...\n", commandList)

		input, err := reader.ReadString('\n')

		if util.IsCancel(err) {
			fmt.Println("Canceled.")
			os.Exit(0)
		}

		input = strings.TrimSpace(input)

		arguments := args.Parse(input)

		if len(arguments) == 0 {
			continue
		}

		command := arguments[0]

		if command == "quit" {
			break
		}

		switch command {
		case "help":
			// pad the top of the help display
			fmt.Println("")

			if len(arguments) == 1 {
				help := gc.Help(&gc.Commands)

				fmt.Println(help)
			} else {
				cmd := strings.Join(arguments[1:], "")

				index := slices.IndexFunc(gc.Commands, func(command gc.Cmd) bool {
					return command.Name == cmd
				})

				if index != -1 {
					help := gc.Commands[index].Help()

					fmt.Println(help)
				} else {
					fmt.Printf("'%v' is not a valid command!\n", cmd)
				}
			}

			fmt.Println("Press `enter` to continue...")

			buf := make([]byte, 1)
			_, err := reader.Read(buf)

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}
		default:
			for _, cmd := range gc.Commands {
				if cmd.Name == command {
					cmd.Run(arguments, game)
					continue
				}
			}

			fmt.Println("Invalid command! Please enter a valid command (catch, kill, end, quit, help)")
		}
	}
}

func delete(game *database.Game) {
	db := database.Connect()
	defer db.Close()

	_, err := db.Exec("DELETE FROM Games WHERE Id = ?", game.Id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted '%v'!\n", game.Name)
}
