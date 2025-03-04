package token

import (
	"DBHS/config"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User interface {
	GetOId() string
	GetUserName() string
}

// takes a user object which implements the User interface and an expiry time for the token
// claims extracted from the user object
func CreateAccessToken(user User, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := TokenClaims{
		Id:       user.GetOId(),      // user oid
		Username: user.GetUserName(), // username
		RegisteredClaims: jwt.RegisteredClaims{ // time of expiry
			ExpiresAt: &jwt.NumericDate{
				Time: time.Unix(exp, 0),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.Env.AccessTokenSecret)
	if err != nil {
		return "", nil
	}
	return accessToken, nil
}

// return a jwt token object if returned an error the token is not valid
func ParseToken(requestToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Env.AccessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// returns nil if the token is a valid token and not expired
func IsAuthorized(requestToken string) error {
	token, err := ParseToken(requestToken)
	if err != nil {
		return err
	}
	// check expiration
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}

	curr := time.Now()
	if curr.After(exp.Time) {
		return fmt.Errorf("Token Expired")
	}

	return nil
}

// return the user id embeded into the token
func GetIdFromToken(requestToken string) (string, error) {
	token, err := ParseToken(requestToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid Token")
	}
	Id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid Token")
	}
	return Id, nil
}

// return the user id embeded into the token
func GetUserNameFromToken(requestToken string) (string, error) {
	token, err := ParseToken(requestToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid Token")
	}
	userName, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid Token")
	}
	return userName, nil
}

func GetData(requestToken string, fields ...string) ([]interface{}, error) {
	token, err := ParseToken(requestToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid Token")
	}

	values := make([]interface{}, len(fields))
	for i, k := range fields {
		v, ok := claims[k]
		if !ok {
			return nil, fmt.Errorf("Invalid Token")
		}
		values[i] = v
	}
	return values, nil
}
