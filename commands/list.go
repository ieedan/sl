package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"database"
	"util"
)

func (l *ListCmd) Run() error {
	db := database.Connect()
	defer db.Close()

	var names []string

	err := db.Select(&names, "Select Name FROM Games G")
	if err != nil {
		log.Fatal(err)
	}

	games := make(map[int64]database.Game)

	for _, name := range names {
		// we know that ok will always be true since we just got it
		game, _ := database.GetGame(db, name)

		games[game.Id] = *game
	}

	slice := util.MapToSlice(&games, func(k int64, v database.Game) database.Game {
		return v
	})

	mappedNames := util.Map(&slice, func(item database.Game, i int) string {
		return item.Name
	})

	mappedNames = append(mappedNames, "Name")

	mappedDates := util.Map(&slice, func(item database.Game, i int) string {
		return item.CreatedAt.Local().Format("01-02-2006 3:04 PM")
	})

	mappedDates = append(mappedDates, "Started")

	mappedTrainers := util.Map(&slice, func(item database.Game, i int) string {
		return strconv.Itoa(len(*item.Trainers))
	})

	mappedTrainers = append(mappedTrainers, "Players")

	nameMin := util.MinLength(&mappedNames)
	trainersMin := util.MinLength(&mappedTrainers)
	dateMin := util.MinLength(&mappedDates)

	table := "\n"

	table += util.LPad(fmt.Sprintf("│ %v │ %v │ %v │\n", util.PadRightMin("Name", nameMin), util.PadRightMin("Players", trainersMin), util.PadRightMin("Started", dateMin)), 2)
	table += util.LPad(fmt.Sprintf("├─%v─┼─%v─┼─%v─┤\n", strings.Repeat("─", nameMin), strings.Repeat("─", trainersMin), strings.Repeat("─", dateMin)), 2)

	for _, game := range games {
		name := util.PadRightMin(game.Name, nameMin)
		trainers := util.PadRightMin(strconv.Itoa(len(*game.Trainers)), trainersMin)
		createdAt := util.PadRightMin(game.CreatedAt.Local().Format("1-2-2006 3:04 PM"), dateMin)

		table += util.LPad(fmt.Sprintf("│ %v │ %v │ %v │\n", name, trainers, createdAt), 2)
	}

	fmt.Println(table)

	return nil
}

type ListCmd struct {
}
