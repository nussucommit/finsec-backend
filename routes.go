package main

func (s *server) routes() {
	// Handler functions don’t actually handle the requests, they return a function that does.
	// This gives us a closure environment in which our handler can operate.
	// If a particular handler has a dependency, take it as an argument.
	// Reference: https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
	s.router.HandleFunc("/", s.handleHelloWord()).Methods("GET")
	s.router.HandleFunc("/user/signup", s.hanldeUserSignUp()).Methods("POST")
}
