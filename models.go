package main

// Define your models here
type response struct {
	Message string `json:"message"`
}
type user struct {
	User_id    int    `json:"user_id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Contact_no string `json:"contact_no"`
	Role       int    `json:"role"`
}

type quotation struct {
	Quotation_id     int    `json:"quotation_id"`
	Event_name       string `json:"event_name"`
	Item_description string `json:"item_description"`
	Item_quantity    int    `json:"item_quantity"`
	Student_name     string `json:"student_name"`
	Status           int    `json:"status"`
}

// type role struct {
// 	Role_id     int    `json:"role_id"`
// 	Role_name   string `json:"role_name"`
// 	Description string `json:"description"`
// }
