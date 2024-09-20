package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Location to open the database from
const Location string = "./database.db"

type Game struct {
	Id        int32
	Name      string
	CreatedAt string
	Trainers  *[]Trainer
}

type Trainer struct {
	Id     int32
	GameId int32
	Name   string
}

type Route struct {
	Id              int32
	GameId          int32
	Name            string
	PokemonAreAlive bool
	Pokemon         *[]Pokemon
}

type Pokemon struct {
	Id        int32
	RouteId   int32
	TrainerId int32
	Name      string
}

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", Location)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetGame(db *sql.DB, name string) Game {
	qry := `
		SELECT 
			Id,
			Name,
			CreatedAt
		FROM Games
		WHERE Name = ?
	`

	stmt, err := db.Prepare(qry)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var game Game

	for rows.Next() {
		rows.Scan(&game.Id, &game.Name, &game.CreatedAt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	trainers := GetTrainers(db, int(game.Id))

	game.Trainers = &trainers

	return game
}

func GetTrainers(db *sql.DB, gameId int) []Trainer {
	qry := `
	SELECT Id, GameId, Name FROM Trainers WHERE GameId = ?
	`

	stmt, err := db.Prepare(qry)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(gameId)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trainers []Trainer

	for rows.Next() {
		var trainer Trainer

		rows.Scan(&trainer.Id, &trainer.GameId, &trainer.Name)

		trainers = append(trainers, trainer)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return trainers
}

func GetRoutes(db *sql.DB, gameId int) []Route {
	qry := `
	SELECT Id, GameId, Name, PokemonAreAlive FROM Routes WHERE GameId = ?
	`

	stmt, err := db.Prepare(qry)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(gameId)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var routes []Route

	for rows.Next() {
		var route Route

		rows.Scan(&route.Id, &route.GameId, &route.Name, &route.PokemonAreAlive)

		routes = append(routes, route)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for i, route := range routes {
		pokemon := GetPokemon(db, int(route.Id))

		routes[i].Pokemon = &pokemon
	}

	return routes
}

func GetPokemon(db *sql.DB, routeId int) []Pokemon {
	qry := `
	SELECT Id, RouteId, TrainerId, Name FROM Pokemon WHERE RouteId = ?
	`

	stmt, err := db.Prepare(qry)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(routeId)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pokemon []Pokemon

	for rows.Next() {
		var p Pokemon

		rows.Scan(&p.Id, &p.RouteId, &p.TrainerId, &p.Name)

		pokemon = append(pokemon, p)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return pokemon
}
