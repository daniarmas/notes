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
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func (ds *hashDatasource) CheckHash(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}