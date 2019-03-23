package client

import (
	"fmt"
	"github.com/moezakura/escape-proxy/model"
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

func (a *Auth) Valid(user, password string) bool {
	if !a.isAuth {
		fmt.Println(">>>")
		return true
	}

	for _, u := range a.users {
		if u.Id == user && u.Password == password {
			return true
		}
	}

	return false
}
