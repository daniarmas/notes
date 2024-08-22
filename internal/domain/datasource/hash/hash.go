package domainhashds

type HashDs interface {
	Hash(value string) (string, error)
	CheckHash(value, hash string) bool
}
