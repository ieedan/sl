package database

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ieedan/sl/table"
	"github.com/ieedan/sl/util"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Path to open the database from
const DatabaseLocation string = "database.db"

const MigrationsLocation string = "migration.sql"

type Game struct {
	Id        int64     `db:"Id"`
	Name      string    `db:"Name"`
	CreatedAt time.Time `db:"CreatedAt"`
	Trainers  *[]Trainer
	Routes    *[]Route
}

func (g *Game) String() string {
	t := table.New(table.DEFAULT_OPTIONS)

	headerColumns := util.Map(g.Trainers, func(trainer Trainer, index int) string {
		return trainer.Name
	})

	headerColumns = append([]string{"Route"}, headerColumns...)

	t.AddHeader(headerColumns...)

	trainerPokemonByRoute := make(map[int64][]string)

	for _, trainer := range *g.Trainers {
		for _, pokemon := range *trainer.Pokemon {
			_, ok := trainerPokemonByRoute[pokemon.RouteId]

			if ok {
				trainerPokemonByRoute[pokemon.RouteId] = append(trainerPokemonByRoute[pokemon.RouteId], pokemon.Name)
				continue
			}

			trainerPokemonByRoute[pokemon.RouteId] = []string{pokemon.Name}
		}
	}

	for _, route := range *g.Routes {
		columns := []string{
			route.Name,
		}

		columns = append(columns, trainerPokemonByRoute[route.Id]...)

		transform := func(str string) string {
			if route.PokemonAreAlive {
				return str
			}

			return util.StrikeThrough(color.RedString(str))
		}

		t.AddRowTransform(transform, columns...)
	}

	return t.String()
}

func (g *Game) GetRoute(name string) (int64, bool) {
	firstTrainer := (*g.Trainers)[0]

	var id int64

	for _, p := range *firstTrainer.Pokemon {
		if strings.EqualFold(p.Route.Name, name) {
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

// Returns a database connection. If the database has not been created it will create it and run migrations
func Connect() *sqlx.DB {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	execPath = filepath.Dir(execPath)

	migrationsPath := filepath.Join(execPath, MigrationsLocation)

	dev := false
	// the principle of this is that the migrations file won't exist in the
	// Go temp directory so we use it to determine our environment
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		dev = true
	}

	// when dev just make the path the place where it was called from
	if dev {
		execPath = "./"
		migrationsPath = filepath.Join(execPath, MigrationsLocation)
	}

	databasePath := filepath.Join(execPath, DatabaseLocation)

	migrate := false
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		migrate = true
	}

	db, err := sqlx.Connect("sqlite", databasePath)

	if err != nil {
		log.Fatal(err)
	}

	// if migrations are necessary automatically run them
	if migrate {
		bytes, err := os.ReadFile(migrationsPath)
		if err != nil {
			log.Fatal(err)
		}

		qry := string(bytes)

		_, err = db.Exec(qry)
		if err != nil {
			log.Fatal(err)
		}
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
	SELECT Id, GameId, Name, PokemonAreAlive FROM Routes WHERE GameId = ? ORDER BY Id ASC
	`

	var routes []Route

	err := db.Select(&routes, qry, gameId)
	if err != nil {
		log.Fatal(err)
	}

	return routes
}
