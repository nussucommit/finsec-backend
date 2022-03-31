package main

import (
	"database/sql"
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

func (s *server) handleUserSignUp() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		User_Id string `json:"user_id"`
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

		var userId string
		row := s.db.QueryRow("SELECT user_id FROM Users WHERE email = $1", newUser.Email)

		if err := row.Scan(&userId); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		resp := response{userId}
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
		var currentUser user
		err := decoder.Decode(&currentUser)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		var password string

		if err = s.db.QueryRow("SELECT password FROM Users WHERE email=$1", currentUser.Email).Scan(&password); err != nil {
			respond(w, r, response{"User not found"}, http.StatusOK)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(currentUser.Password)); err != nil {
			respond(w, r, response{"Incorrect password"}, http.StatusOK)
			return
		}

		var token string

		if token, err = generateToken(currentUser); err != nil {
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

func (s *server) handleResetPassword() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Email  string `json:"email"`
		UserId string `json:"userId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		// var data map[string]interface{}
		var usernameOrEmail string

		err := decoder.Decode(&usernameOrEmail)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		var email string
		var userId string
		row := s.db.QueryRow("SELECT email, user_id FROM Users WHERE name = $1", usernameOrEmail)

		if err := row.Scan(&email, &userId); err != nil {
			if err == sql.ErrNoRows {
				row = s.db.QueryRow("SELECT email, user_id FROM Users WHERE email = $1", usernameOrEmail)
				if err = row.Scan(&email, &userId); err != nil {
					respondErr(w, r, err, http.StatusInternalServerError)
					return
				}
			}
		}

		resp := response{email, userId}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleResetPasswordEmail() func(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email string `json:"email"`
		URL   string `json:"url"`
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		// var data map[string]interface{}
		var emailURL request

		err := decoder.Decode(&emailURL)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		setResetPasswordEmail([]string{emailURL.Email}, emailURL.URL)

		resp := response{"Email sent successfully"}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleChangePassword() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		// var data map[string]interface{}
		var newPassword string
		var userId = r.URL.Query()["userId"]
		err := decoder.Decode(&newPassword)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		passEncrypted, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

		_, err = s.db.Exec("INSERT INTO Users (password) VALUES $1 WHERE user_id = $2",
			passEncrypted,
			userId,
		)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		resp := response{"Password Changed"}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleSignupConfirmationEmail() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decoder := json.NewDecoder(r.Body)

		var userId = r.URL.Query()["userId"]
		row := s.db.QueryRow("SELECT email FROM Users WHERE user_id = $1", userId)

		var email string

		if err := row.Scan(&email); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		confirmSignUpEmail([]string{email})

		resp := response{"Confirmation email sent"}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleVerifiedAccountStatus() func(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decoder := json.NewDecoder(r.Body)

		var userId = r.URL.Query()["userId"]
		row := s.db.QueryRow("SELECT email, account_status FROM Users WHERE user_id = $1", userId)

		var email string
		var account_status string

		if err := row.Scan(&email, &account_status); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		if account_status == "verified" {
			resp := response{"Account status already verified"}
			respond(w, r, resp, http.StatusBadRequest)
			return
		}

		_, err := s.db.Exec("UPDATE Users SET account_status = $1 WHERE user_id = $2",
			"verified",
			userId,
		)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		resp := response{"Account status verified"}
		respond(w, r, resp, http.StatusOK)
	}
}
