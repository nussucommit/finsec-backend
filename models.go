package main

// Define your models here
type user struct {
	User_id    int    `json:"user_id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Salt       string `json:"salt"`
	Contact_no string `json:"contact_no"`
	Role       int    `json:"role"`
}

type role struct {
	Role_id     int    `json:"role_id"`
	Role_name   string `json:"role_name"`
	Description string `json:"description"`
}
