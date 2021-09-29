package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rrr_goBackend/models"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	Id:       10,
	Email:    "me@here.com",
	Password: "$2a$12$86jCqVSqV3EMwFWbpY1OoeRV1hiVYLFb9C3CUXGs3Z0zTdmvzXvY.",
}

var unauthErr = errors.New("unauthorized")
var jwtErr = errors.New("error signing")

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Calling signin.")
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, unauthErr)
		return
	}

	hashedPW := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, unauthErr)
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.Id)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "myDomain.com"
	claims.Audiences = []string{"myDomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, jwtErr)
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}
