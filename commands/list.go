package commands

import (
	"fmt"
	"log"
	"strconv"

	"database"
	"util"
)

func (l *ListCmd) Run() error {
	db := database.Connect()
	defer db.Close()

	qry := `
		Select G.* FROM Games G
	`

	rows, err := db.Query(qry)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	games := make(map[int32]database.Game)

	for rows.Next() {
		var game database.Game

		err = rows.Scan(&game.Id, &game.Name, &game.CreatedAt)

		if _, ok := games[game.Id]; !ok {
			games[game.Id] = game
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	slice := util.MapToSlice(&games, func(k int32, v database.Game) database.Game {
		return v
	})

	mappedNames := util.Map(&slice, func(item database.Game, i int) string {
		return item.Name
	})

	mappedDates := util.Map(&slice, func(item database.Game, i int) string {
		return item.CreatedAt
	})

	mappedIds := util.Map(&slice, func(item database.Game, i int) string {
		return strconv.Itoa(int(item.Id))
	})

	idMin := util.MinLength(&mappedIds)
	nameMin := util.MinLength(&mappedNames)
	dateMin := util.MinLength(&mappedDates)

	for _, game := range games {
		id := util.PadLeftMin(strconv.Itoa(int(game.Id)), int32(idMin))
		name := util.PadRightMin(game.Name, int32(nameMin))
		createdAt := util.PadRightMin(game.CreatedAt, int32(dateMin))

		fmt.Printf("%v | %v | %v\n", id, name, createdAt)
	}

	return nil
}

type ListCmd struct {
}
