package models

type RegisterModelReq struct {
	Id          string `json:"-"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type AuthorizationReq struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type Authorization struct {
	Token  string `json:"token"`
	Status bool   `json:"status"`
}

type RegisterModelRes struct {
	Status bool `json:"status"`
}

type RegisterRes struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	AccessToken string `json:"access_token"`
}
