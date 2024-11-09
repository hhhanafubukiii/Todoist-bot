package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type postgres struct {
	db *pgxpool.Pool
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (pg *postgres) SaveAccessToken(ctx context.Context, chatId int64, accessToken string) error {
	conn := getConn()
	query := fmt.Sprintf(`INSERT INTO tokens(chat_id, access_token) VALUES (%d, %s)`, chatId, accessToken)
	_, err := conn.Exec(ctx, query, chatId, accessToken)
	if err != nil {
		log.Fatal("unable to insert access token!", err)
	}

	return nil
}

func (pg *postgres) GetAccessToken(ctx context.Context, chatId int64) (string, error) {
	conn := getConn()
	defer conn.Close(context.Background())
	query := fmt.Sprintf(`SELECT access_token FROM tokens WHERE chat_id = %d`, chatId)

	row, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Fatal("unable to get access token!", err)
	}
	defer row.Close()

	var accessToken string

	err = row.Scan(&accessToken)
	if err != nil {
		log.Fatal("unable to get access token!", err)
	}

	return accessToken, nil
}

func getConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("databaseURL"))
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	return conn
}
