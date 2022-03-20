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

func (s *server) handleUserSignUp() func(w http.ResponseWriter, r *http.Request) {
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

func (s *server) handleQuotationGetAll() func(w http.ResponseWriter, r *http.Request) {

	type responseQuotation struct {
		Quotations []quotation `json:"quotattions"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var quotations []quotation

		sqlQuery := `SELECT quotation_id, 
					event_name, 
					item_description, 
					item_quantity, 
					student_name, 
					status 
					FROM Quotations`

		rows, err := s.db.Query(sqlQuery)

		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var qtn quotation

			err = rows.Scan(
				&qtn.Quotation_id,
				&qtn.Event_name,
				&qtn.Item_description,
				&qtn.Item_quantity,
				&qtn.Student_name,
				&qtn.Status,
			)

			if err != nil {
				respondErr(w, r, err, http.StatusInternalServerError)
				return
			}

			quotations = append(quotations, qtn)
		}

		responseQuotation := responseQuotation{Quotations: quotations}

		respond(w, r, responseQuotation, http.StatusOK)

	}
}

// Handler function to change status to approved/rejected
// Resubmit quotation could also be done by changing the status (?)
func (s *server) handleQuotationUpdateStatus() func(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Quotation_id int `json:"quotation_id"`
		New_status   int `json:"new_status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var quotationRequest request

		if err := decoder.Decode(&quotationRequest); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
		}

		sqlQuery := `UPDATE Quotations 
					SET status = $1
					WHERE quotation_id = $2`

		if _, err := s.db.Exec(sqlQuery, quotationRequest.New_status, quotationRequest.Quotation_id); err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		response := response{"Status has been changed"}
		respond(w, r, response, http.StatusOK)
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
