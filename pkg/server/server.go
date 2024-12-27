package server

import (
	"Todoist-bot/pkg/storage/postgres"
	"context"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type AuthorizationServer struct {
	server        *http.Server
	todoistClient *todoist.Client
	db            *postgres.Postgres
	URL           string
}

var (
	clientID     = os.Getenv("client_id")
	clientSecret = os.Getenv("client_secret")
	databaseURL  = os.Getenv("databaseURL")
)

func NewAuthorizationServer(db *postgres.Postgres, todoistClient *todoist.Client, URL string) *AuthorizationServer {
	return &AuthorizationServer{
		todoistClient: todoistClient,
		db:            db,
		URL:           URL,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	codeParam := r.URL.Query().Get("code")
	stateParam := r.URL.Query().Get("state")
	if codeParam == "" || stateParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, _ := strconv.ParseInt(stateParam, 10, 64)

	accessToken, err := s.todoistClient.GetAccessToken(clientID, clientSecret, codeParam)
	if err != nil {
		log.Fatal("cannot get access token: ", err)
	}

	err = s.db.Save(context.Background(), chatID, accessToken, databaseURL)
	if err != nil {
		log.Fatal("cannot save access token: ", err)
	}
}
