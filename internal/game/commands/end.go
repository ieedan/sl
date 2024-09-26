package game

import (
	"github.com/ieedan/sl/internal/database"
	"github.com/ieedan/sl/internal/util"
)

var End = Cmd{Name: "end", Description: "Ends the game (whiteout / blackout)", Run: end}

func end(args []string, game *database.Game) {
	ids := util.Map(&game.Routes, func(route database.Route, i int) int64 {
		return route.Id
	})

	KillRoutes(ids...)
}
