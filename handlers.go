package main

import (
	"net/http"
)

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
