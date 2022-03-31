package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func sendEmail(recipients []string, subject string, message string) {
	from := ""     // commIT email here
	password := "" // commIT password here
	host := "smtp.gmail.com"

	// Its the default port of smtp server
	port := "587"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := []byte("Subject: " + subject + "\r\n" + mime + message)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, recipients, body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully sent mail to all user in toList")
}

func setResetPasswordEmail(recipients []string, url string) {
	subject := "Reset Password"

	// use html so can add link/logo etc
	message :=
		`
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<title>Reset Password</title>
		<style>
			body {
				background-color: #FFFFFF; padding: 0; margin: 0;
			}
		</style>
	</head>

	<body>
		<span style="font-size: 12px; line-height: 1.5; color: #333333;">
			We have received a request to reset your password.
			<br/><br/>

			To reset your password, please follow the link below:
			<br/>
			<a href=` + url + `>reset password link</a>

			<br/><br/>

			If you did not request for a change in password, please contact us immediately.
		</span>
	</body>
	`
	sendEmail(recipients, subject, message)
}

func confirmSignUpEmail(recipients []string) {
	subject := "Confirm Registration"

	// use html so can add link/logo etc
	message :=
		`
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<title>Confirm Registration</title>
		<style>
			body {
				background-color: #FFFFFF; padding: 0; margin: 0;
			}
		</style>
	</head>

	<body">
		<span style="font-size: 12px; line-height: 1.5; color: #333333;">
			This email is to acknowledge that your registration under this email address is successful.
			<br/><br/>

			If you did not sign up under this email, please contact us immediately.
		</span>
	</body>
	`
	sendEmail(recipients, subject, message)

}
