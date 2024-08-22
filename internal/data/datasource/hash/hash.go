package datahashds

import (
	"github.com/daniarmas/notes/internal/domain/datasource/hash"
	"golang.org/x/crypto/bcrypt"
)

type hashDatasource struct {
}

func NewBcryptHashDatasource() domainhashds.HashDs {
	return &hashDatasource{}
}

func (ds hashDatasource) Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func (ds hashDatasource) CheckHash(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}
