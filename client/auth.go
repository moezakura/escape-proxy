package main

import "github.com/moezakura/escape-proxy/model"

type Auth struct {
	users []model.AuthUsers
}

func NewAuth(users []model.AuthUsers) Auth {
	return Auth{
		users: users,
	}
}

func (a Auth) Valid(user, password string) bool {
	for _, u := range a.users{
		if u.Id == user && u.Password == password{
			return true
		}
	}

	return false
}
