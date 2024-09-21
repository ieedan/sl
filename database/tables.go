package database

import (
	"fmt"
	"log"
	"strings"
	"util"

	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Location to open the database from
const Location string = "./database.db"

type Game struct {
	Id        int64  `db:"Id"`
	Name      string `db:"Name"`
	CreatedAt string `db:"CreatedAt"`
	Trainers  *[]Trainer
	Routes    *[]Route
}

func (g *Game) String() string {
	routes := []string{
		"Route",
	}

	firstTrainer := (*g.Trainers)[0]

	for _, pokemon := range *firstTrainer.Pokemon {
		routes = append(routes, pokemon.Route.Name)
	}

	routesMin := util.MinLength(&routes)

	trainerMinByIndex := make(map[int]int)

	for i, trainer := range *g.Trainers {
		pokemon := []string{
			trainer.Name,
		}

		for _, p := range *trainer.Pokemon {
			pokemon = append(pokemon, p.Name)
		}

		trainerMinByIndex[i] = util.MinLength(&pokemon)
	}

	var table string = "\n"

	for i, route := range routes {
		var row string = "│"
		row += fmt.Sprintf(" %v │", util.PadRightMin(route, routesMin))

		isDead := false

		for trainerIndex, trainer := range *g.Trainers {
			min := trainerMinByIndex[trainerIndex]
			// for the heading
			if i == 0 {
				row += fmt.Sprintf(" %v │", util.PadRightMin(trainer.Name, min))
				continue
			}

			pokemon := (*trainer.Pokemon)[i-1]

			isDead = !pokemon.Route.PokemonAreAlive

			row += fmt.Sprintf(" %v │", util.PadRightMin(pokemon.Name, min))
		}

		if isDead {
			row = util.StrikeThrough(color.RedString(row))
		}

		table += util.LPad(row, 2) + "\n"

		if i == 0 {
			row = "├"

			row += strings.Repeat("─", routesMin+2) + "┼"

			for trainerIndex := range *g.Trainers {
				min := trainerMinByIndex[trainerIndex]

				row += strings.Repeat("─", min+2)

				if trainerIndex+1 < len(*g.Trainers) {
					row += "┼"
				} else {
					row += "┤"
				}
			}

			table += util.LPad(row, 2) + "\n"
		}
	}

	return table
}

func (g *Game) GetRoute(name string) (int64, bool) {
	firstTrainer := (*g.Trainers)[0]

	var id int64

	for _, p := range *firstTrainer.Pokemon {
		if p.Route.Name == name {
			id = p.Route.Id
		}
	}

	return id, id != 0
}

func (g *Game) IsDead() bool {
	for _, route := range *g.Routes {
		if route.PokemonAreAlive {
			return false
		}	
	}

	return true
}

type Trainer struct {
	Id      int64  `db:"Id"`
	GameId  int64  `db:"GameId"`
	Name    string `db:"Name"`
	Pokemon *[]Pokemon
}

type Route struct {
	Id              int64  `db:"Id"`
	GameId          int64  `db:"GameId"`
	Name            string `db:"Name"`
	PokemonAreAlive bool   `db:"PokemonAreAlive"`
}

type Pokemon struct {
	Id        int64  `db:"Id"`
	RouteId   int64  `db:"RouteId"`
	TrainerId int64  `db:"TrainerId"`
	Name      string `db:"Name"`
	Route     *Route
}

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", Location)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetGame(db *sqlx.DB, name string) (*Game, bool) {
	qry := `
		SELECT 
			Id,
			Name,
			CreatedAt
		FROM Games
		WHERE Name = ?
	`
	var game Game

	err := db.QueryRowx(qry, name).StructScan(&game)
	if err != nil {
		return nil, false
	}

	trainers := GetTrainers(db, game.Id)

	game.Trainers = &trainers

	routes := GetRoutes(db, game.Id)

	game.Routes = &routes

	return &game, true
}

func GetTrainers(db *sqlx.DB, gameId int64) []Trainer {
	qry := `
	SELECT Id, GameId, Name FROM Trainers WHERE GameId = ?
	`
	var trainers []Trainer

	err := db.Select(&trainers, qry, gameId)
	if err != nil {
		log.Fatal(err)
	}

	for i, trainer := range trainers {
		pokemon := GetPokemon(db, trainer.Id)

		trainers[i].Pokemon = &pokemon
	}

	return trainers
}

func GetPokemon(db *sqlx.DB, trainerId int64) []Pokemon {
	qry := `
	SELECT Id, RouteId, TrainerId, Name FROM Pokemon WHERE TrainerId = ? ORDER BY RouteId ASC
	`

	var pokemon []Pokemon

	err := db.Select(&pokemon, qry, trainerId)
	if err != nil {
		log.Fatal(err)
	}

	for i, p := range pokemon {
		route := GetRoute(db, p.RouteId)

		pokemon[i].Route = &route
	}

	return pokemon
}

func GetRoute(db *sqlx.DB, routeId int64) Route {
	qry := `
	SELECT Id, GameId, Name, PokemonAreAlive FROM Routes WHERE Id = ?
	`

	var route Route

	err := db.Get(&route, qry, routeId)
	if err != nil {
		log.Fatal(err)
	}

	return route
}

func GetRoutes(db *sqlx.DB, gameId int64) []Route {
	qry := `
	SELECT Id, GameId, Name, PokemonAreAlive FROM Routes WHERE GameId = ?
	`

	var routes []Route

	err := db.Select(&routes, qry, gameId)
	if err != nil {
		log.Fatal(err)
	}

	return routes
}
