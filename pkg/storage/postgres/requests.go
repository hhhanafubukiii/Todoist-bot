package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

type Postgres struct {
	db *pgxpool.Pool
}

func (pg *Postgres) SaveAccessToken(ctx context.Context, chatId int64, accessToken string, databaseURL string) error {
	conn := getConn(databaseURL)
	query := fmt.Sprintf(`INSERT INTO tokens(chat_id, access_token) VALUES (%d, %s)`, chatId, accessToken)
	_, err := conn.Exec(ctx, query, chatId, accessToken)
	if err != nil {
		log.Fatal("unable to insert access token!", err)
	}

	return nil
}

func (pg *Postgres) GetAccessToken(ctx context.Context, chatId int64, databaseURL string) (accessToken string, err error) {
	conn := getConn(databaseURL)
	query := fmt.Sprintf(`SELECT access_token FROM tokens WHERE chat_id = %d`, chatId)

	row, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Fatal("unable to get access token!", err)
	}

	err = row.Scan(&accessToken)
	if err != nil {
		log.Fatal("unable to get access token!", err)
	}

	defer conn.Close(context.Background())
	defer row.Close()

	return
}

func getConn(databaseURL string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	return conn
}
