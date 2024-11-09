package server

import "net/http"

type AuthorizationServer struct {
	server *http.Server
	// todoistClient *todoist.Client
	// tokenRepository repository.TokenRepository
	redirectURL string
}

// func NewAuthorizationServer (todoistClient *todoist.Client, tokenRepository repository.TokenRepository,redirectURL string) *AuthorizationServer {
// return &AuthorizationServer{todoistClient: todoistClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
// }

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
