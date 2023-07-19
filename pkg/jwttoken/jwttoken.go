package jwttoken

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func AddJwtToken(w http.ResponseWriter, user string) error {
	//mySigningKey := []byte(s.db.Select(s.user.User)) //hash
	mySigningKey := []byte("SecretYouShouldHide")
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(mySigningKey)

	if err != nil {
		return err
	}
	w.Header().Add("Jwt-Token", tokenStr)
	return nil
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
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("SecretYouShouldHide"), nil
	})

	if err != nil {
		//fmt.Printf("step1: claim data [%v]", claims)
		return errors.New("Can't parse JWT token."), false
	}

	if token == nil || !token.Valid {
		//fmt.Printf("step2: %v\n", token.Valid)
		return errors.New("Can't parse JWT token."), false
	}

	//fmt.Printf("[%v]\n", token)

	return nil, true
}
