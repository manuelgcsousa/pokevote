package data

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DbDirPath string = "~/.local/share/pokevote"

const PokemonTableSchema string = `
CREATE TABLE IF NOT EXISTS pokemon (
    id          INT PRIMARY KEY,
    dex_id      INT,
    name        VARCHAR(255),
    sprite_url  VARCHAR(255),
    votes_won   INT DEFAULT 0,
    votes_lost  INT DEFAULT 0
);
`

type Pokemon struct {
	Id        int    `json:"id"         db:"id"`
	DexId     int    `json:"dex_id"     db:"dex_id"`
	Name      string `json:"name"       db:"name"`
	SpriteUrl string `json:"sprite_url" db:"sprite_url"`
	VotesWon  int    `json:"votes_won"  db:"votes_won"`
	VotesLost int    `json:"votes_lost" db:"votes_lost"`
}

type Database struct {
	conn *sqlx.DB
}

func NewConnection() (*Database, error) {
	dbDirPath := DbDirPath

	// Get user's HOME directory path
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
	dbDirPath = filepath.Join(home, dbDirPath[2:])

	// Ensure DB dir exists
	if err := os.MkdirAll(dbDirPath, 0755); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}

	// Connect to DB
	dbSrcName := fmt.Sprintf("%s%s", dbDirPath, "/poke.db")
	db, err := sqlx.Connect("sqlite3", dbSrcName)
	if err != nil {
		return nil, err
	}

	return &Database{conn: db}, nil
}

func (db *Database) InitDatabase() {
	db.conn.MustExec(PokemonTableSchema)
}

func (db *Database) CloseDatabase() {
	db.conn.Close()
}

func (db *Database) InsertPokemon(pokemon *Pokemon) error {
	query := `INSERT INTO pokemon (id, dex_id, name, sprite_url) VALUES (:id, :dex_id, :name, :sprite_url)`

	_, err := db.conn.NamedExec(query, pokemon)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTwoRandomPokemon() *[]Pokemon {
	query := "SELECT * FROM pokemon WHERE sprite_url IS NOT NULL AND sprite_url != '' ORDER BY RANDOM() LIMIT 2"
	pokemon := []Pokemon{}

	err := db.conn.Select(&pokemon, query)
	if err != nil {
		slog.Default().Error("Error (GetTwoRandomPokemon): " + err.Error())
	}

	return &pokemon
}

func (db *Database) GetAllPokemon() *[]Pokemon {
	query := `SELECT * FROM pokemon
	WHERE (sprite_url IS NOT NULL AND sprite_url != '')
	ORDER BY votes_won DESC, votes_lost ASC`

	pokemon := []Pokemon{}

	err := db.conn.Select(&pokemon, query)
	if err != nil {
		slog.Default().Error("Error (GetAllPokemon): " + err.Error())
	}

	return &pokemon
}

func (db *Database) UpdatePokemonVoteCounters(idWithVoteWon int, idWithVoteLost int) {
	query_votes_won := `UPDATE pokemon SET votes_won = votes_won + 1 WHERE id = ?`
	query_votes_lost := `UPDATE pokemon SET votes_lost = votes_lost + 1 WHERE id = ?`

	// Begin transaction
	tx, err := db.conn.Begin()
	if err != nil {
		slog.Default().Error("Error while starting DB transaction: " + err.Error())
		return
	}

	// Update votes won counter
	_, err = tx.Exec(query_votes_won, idWithVoteWon)
	if err != nil {
		slog.Default().Error("Error while updating 'votes_won' counter: " + err.Error())
	}

	// Update votes lost counter
	_, err = tx.Exec(query_votes_lost, idWithVoteLost)
	if err != nil {
		slog.Default().Error("Error while updating 'votes_lost' counter: " + err.Error())
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		slog.Default().Error("Error while commiting transaction: " + err.Error())
	}
}
