package models

type ClientResponse struct {
	Client ClientDetails `json:"client"`
}

type ClientDetails struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Company   string `json:"companyname"`
	Address   string `json:"address1"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phonenumber"`
}
