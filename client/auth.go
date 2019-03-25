package client

import (
	"github.com/moezakura/escape-proxy/model"
	"net"
)

type Auth struct {
	users  []model.AuthUsers
	isAuth bool
}

func NewAuth(isAuth bool, users []model.AuthUsers) *Auth {
	return &Auth{
		users:  users,
		isAuth: isAuth,
	}
}

func (a *Auth) Authenticate(user, password string, addr net.Addr) bool {
	if !a.isAuth {
		return true
	}

	for _, u := range a.users {
		if u.Id == user && u.Password == password {
			return true
		}
	}

	return false
}
