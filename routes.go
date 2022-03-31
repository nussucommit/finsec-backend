package main

func (s *server) routes() {
	// Handler functions don’t actually handle the requests, they return a function that does.
	// This gives us a closure environment in which our handler can operate.
	// If a particular handler has a dependency, take it as an argument.
	// Reference: https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
	s.router.HandleFunc("/", s.handleHelloWord()).Methods("GET")
	s.router.HandleFunc("/user/signup", s.handleUserSignUp()).Methods("POST")
	s.router.HandleFunc("/user/signin", s.handleUserSignIn()).Methods("POST")
	s.router.HandleFunc("/user/reset_password", s.handleResetPassword()).Methods("POST")
	s.router.HandleFunc("/user/reset_password_email", s.handleResetPasswordEmail()).Methods("POST")
	s.router.HandleFunc("/user/change_password/:userId", s.handleChangePassword()).Methods("POST")
	s.router.HandleFunc("/user/signup_confirmation_email/:userId", s.handleSignupConfirmationEmail()).Methods("POST")
	s.router.HandleFunc("/user/verified_account_status/:userId", s.handleVerifiedAccountStatus()).Methods("POST")
}
