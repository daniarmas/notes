package data

import (
	"github.com/daniarmas/notes/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type hashDatasource struct {
}

func NewBcryptHashDatasource() domain.HashDatasource {
	return &hashDatasource{}
}

func (ds *hashDatasource) Hash(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedValue), nil
}

func (ds *hashDatasource) CheckHash(value, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	switch err {
	case nil:
		return true, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	default:
		return false, err
	}
}
