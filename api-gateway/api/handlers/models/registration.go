package models

type RegisterModelReq struct {
	Id        string `json:"-"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Bio       string `json:"bio"`
	Website   string `json:"website"`
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
