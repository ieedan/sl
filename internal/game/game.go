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
	"github.com/ieedan/sl/internal/util"

	tm "github.com/buger/goterm"
	"github.com/jmoiron/sqlx"
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
			fmt.Println("Waiting for command (catch, kill, end, quit, help)...")

			input, err := reader.ReadString('\n')

			if util.IsCancel(err) {
				fmt.Println("Canceled.")
				os.Exit(0)
			}

			input = strings.TrimSpace(input)

			if input == "" {
				continue
			}

			endCommandIndex := strings.Index(input, " ")

			var command string
			args := ""

			if endCommandIndex == -1 {
				command = input
			} else {
				command = input[0:endCommandIndex]
				args = input[endCommandIndex+1:]
			}

			if command == "quit" {
				break
			}

			switch command {
			case "kill":
				var routeName string

				if args != "" {
					routeName = args
				}

				kill(game, &routeName)
			case "catch":
				catch(game)
			case "end":
				end(game)
			case "help":
				fmt.Println("")

				if args == "" {
					help := Help(&Commands)

					fmt.Println(help)
				} else {
					cmd := args

					index := slices.IndexFunc(Commands, func(command Cmd) bool {
						return command.Name == cmd
					})

					if index != -1 {
						help := Commands[index].Help()

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
				fmt.Println("Invalid command! Please enter a valid command (catch, kill, end, quit, help)")
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

func end(game *database.Game) {
	ids := util.Map(&game.Routes, func(route database.Route, i int) int64 {
		return route.Id
	})

	killRoutes(ids...)
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

func kill(game *database.Game, nameOfRoute *string) {
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

func catch(game *database.Game) {
	
}
