package domain

type HashDatasource interface {
	Hash(value string) (string, error)
	CheckHash(value, hash string) bool
}
