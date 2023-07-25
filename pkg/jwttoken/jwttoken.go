package jwttoken

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// Regular
type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func AddJwtToken(user string) (string, error) {
	//mySigningKey := []byte(s.db.Select(s.user.User)) //hash
	mySigningKey := []byte("SecretYouShouldHide")
	// Create the JWT claims, which includes the username and expiry time
	if user == "" {
		return "", errors.New("Claim data incorrect")
	}

	//For tests old
	/*claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}*/
	//For tests old
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//For tests new
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Minute * 30).Unix()})
	tokenStr, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidateJwtToken(r *http.Request) (error, bool) {
	// Get the JWT string from the cookie
	if r.Header["Jwt-Token"] == nil {
		//s.writeResponseError(w, errors.New("Cant find JWT in incoming data.").Error(), http.StatusBadRequest)
		return errors.New("Cant find JWT in incoming data."), false
	}

	tokenStr := r.Header.Get("Jwt-Token")
	//Claims
	claims := &Claims{}
	//Parse token

	//For tests mem.*
	/*token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("SecretYouShouldHide"), nil
	})*/

	//For tests mem.*
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("SecretYouShouldHide"), nil
	})

	if err != nil {
		return err, false
	}

	if token == nil || !token.Valid {
		return errors.New("Can't parse JWT token."), false
	}

	//fmt.Printf("[%v]\n", token)

	return nil, true
}
