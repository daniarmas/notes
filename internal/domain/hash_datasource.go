package domain

// HashDatasource defines the methods for hashing values and checking hashes.
// Implementations of this interface should provide mechanisms to hash values
// and verify if a given hash matches a hashed value.
type HashDatasource interface {
	// Hash takes a string value and returns its hash representation.
	// It returns an error if the hashing operation fails.
	Hash(value string) (string, error)
	// CheckHash verifies if the given hash matches the hash of the provided value.
	// It returns true if the hash matches, false otherwise.
	CheckHash(value, hash string) (bool, error)
}
