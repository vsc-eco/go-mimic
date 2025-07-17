package services

import "errors"

func CreateUser(account, password string) error {
	if account == "" || password == "" {
		return errors.New("invalid account or password")
	}
	return nil
}
