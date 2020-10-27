package main

import (
	"context"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
)

var (
	ab           = authboss.New()
	database     = NewMemStorer()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func SetupAuthboss() *authboss.Authboss {
	ab.Config.Paths.RootURL = "http://localhost:3000"

	ab.Config.Modules.LogoutMethod = "GET"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	defaults.SetCore(&ab.Config, true, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}

	return ab
}

type User struct {
	ID       int
	Email    string
	Password string
}

func (u User) GetPID() string     { return u.Email }

func (u *User) PutPID(pid string) { u.Email = pid }

type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"rick@councilofricks.com": {
				ID:       1,
				Password: "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:    "rick@councilofricks.com",
			},
		},
		Tokens: make(map[string][]string),
	}
}

func (m MemStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)
	m.Users[u.Email] = *u
	return nil
}

func (m MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {
	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}
	return &u, nil
}
