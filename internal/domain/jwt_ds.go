package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/notes/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtDatasource interface {
	CreateJWT(tokenMetadata *JWTMetadata, expirationTime time.Time) error
	ParseJWT(tokenMetadata *JWTMetadata) error
}

type JWTMetadata struct {
	TokenId uuid.UUID
	Token   string
}

type jwtTokenDatasource struct {
	Config *config.Configuration
}

func NewJWTDatasource(cfg *config.Configuration) JwtDatasource {
	return &jwtTokenDatasource{
		Config: cfg,
	}
}

func (j *jwtTokenDatasource) CreateJWT(tokenMetadata *JWTMetadata, expirationTime time.Time) error {
	hmacSecret := []byte(j.Config.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   tokenMetadata.TokenId.String(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return err
	}
	tokenMetadata.Token = tokenString
	return nil
}

func (r *jwtTokenDatasource) ParseJWT(tokenMetadata *JWTMetadata) error {
	hmacSecret := []byte(r.Config.JwtSecret)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenMetadata.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tokenErr := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(tokenErr)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := fmt.Sprintf("%s", claims["sub"])
		tokenIdValue := uuid.MustParse(data)
		tokenMetadata.TokenId = tokenIdValue
		return nil
	} else {
		return err
	}
}