package utils

import (
	"fmt"
	"DBHS/config"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	token 		*jwt.Token
	claims		jwt.MapClaims
}

// returns a new empty token signed with HS256
// header contains the typ set to "JWT" and alg set to "HS256"
func NewToken() Token {
	t := Token{
		token:	jwt.New(jwt.SigningMethodHS256),
	}
	return t
}

// returns new token after parsing the gaven string token 
// tokenString is the raw token inside the header in the HTTP request
func NewTokenString(tokenString string) (Token, error) {
	var token Token
	jwtToken, err := jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) { // returns the secret key to be used for verfiacation
		if config.Secret_Key == "" {
			return nil, fmt.Errorf("could not get the secret Key from the env variables")
		}

		return []byte(config.Secret_Key), nil
	})
	// check if the jwtToken is valid
	if err != nil || !jwtToken.Valid{
		return token, err
	}

	if cliams, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		token = Token{
			token: jwtToken,
			claims: cliams,
		}
	}

	return token, nil
}

// set the header to the passed value 
// this will override the existing value if the header exists
func (t Token) AddHeader(header string, value interface{}) {
	t.token.Header[header] = value
}
// add a set of headers to the token
// this will override the existing value if the header exists
func (t Token) AddHeaders(headers map[string]interface{}) {
	for k, v := range headers {
		t.token.Header[k] = v
	}
}
// return the value of the header if exists and a boolen to indecate if the header exist
func (t Token) GetHeader(header string) (interface{}, bool) {
	value, ok := t.token.Header[header] 
	return value, ok
}

// set the claim to the passed value 
// this will override the existing value if the claim exists
func (t Token) AddClaim(claim string, value interface{}) {
	t.claims[claim] = value
}
// add a set of claims to the token
// this will override the existing value if the claim exists
func (t Token) AddClaims(claims map[string]interface{}) {
	for k, v := range claims {
		t.claims[k] = v
	}
}

// return the value of the claim if exists and a boolen to indecate if the claim exist
func (t Token) GetClaim(claim string) (interface{}, bool) {
	value, ok := t.claims[claim]
	return value, ok
}

func (t Token) String() (string, error) {
	// join the claims
	t.token.Claims = t.claims
	return t.token.SignedString([]byte(config.Secret_Key))
}
