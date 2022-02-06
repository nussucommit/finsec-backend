package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// All handlers are methods of our server struct. So this is how they can access the database functions in db.go (since they are also methods of the server struct)

// For reference only
func (s *server) handleHelloWord() http.HandlerFunc {
	// Define the structure of your request body (for POST requests) and response

	//type request struct {
	//
	//}
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//decoder := json.NewDecoder(r.Body)

		//var req request
		//err := decoder.Decode(&req)
		//if err != nil {
		//	respondErr(w, r, err, http.StatusInternalServerError)
		//	return
		//}
		// Now you can use the request body

		// Do something

		res := response{"Hello World!"}

		respond(w, r, res, http.StatusOK)
	}
}

func (s *server) hanldeUserSignUp() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		// var data map[string]interface{}
		var newUser user
		err := decoder.Decode(&newUser)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		passEncrypted, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)

		_, err = s.db.Exec("INSERT INTO Users (name, password, email) VALUES ($1, $2, $3)",
			newUser.Name,
			passEncrypted,
			newUser.Email,
		)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		resp := response{"User registered"}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleUserSignIn() func(w http.ResponseWriter, r *http.Request) {

	type responseToken struct {
		Message string `json:"token"`
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		// var data map[string]interface{}
		var newUser user
		err := decoder.Decode(&newUser)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		var password string

		if err = s.db.QueryRow("SELECT password FROM Users WHERE email=$1", newUser.Email).Scan(&password); err != nil {
			respond(w, r, response{"User not found"}, http.StatusOK)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(newUser.Password)); err != nil {
			respond(w, r, response{"Incorrect password"}, http.StatusOK)
			return
		}

		var token string

		if token, err = generateToken(newUser); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		resp := responseToken{token}
		respond(w, r, resp, http.StatusOK)
	}
}

func generateToken(u user) (string, error) {
	atClaim := jwt.MapClaims{}
	atClaim["authorized"] = true
	atClaim["user_email"] = u.Email
	atClaim["exp"] = time.Now().Add(time.Minute * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaim)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", err
	}
	return token, nil
}
