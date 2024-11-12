package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Postgres struct {
	db *sql.DB
}

var databaseURL string = os.Getenv("dbURL")

func (pg *Postgres) SaveAccessToken(ctx context.Context, chatId int64, accessToken string, databaseURL string) error {
	db := getConn(databaseURL)
	query := fmt.Sprintf(`INSERT INTO tokens(chat_id, access_token) VALUES (%d, '%s')`, chatId, accessToken)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("unable to insert access token!", err)
	}

	return nil
}

func (pg *Postgres) GetAccessToken(chatId int64, databaseURL string) (accessToken string, err error) {
	db := getConn(databaseURL)
	defer db.Close()
	query := fmt.Sprintf(`SELECT access_token FROM tokens WHERE chat_id = %d`, chatId)

	err = db.QueryRow(query).Scan(&accessToken)
	if err != nil {
		return "", err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return
}

func getConn(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	return db
}
