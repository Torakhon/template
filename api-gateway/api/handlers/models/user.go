package models

type User struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
}

type CreateUserReq struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
}

type GetUserReq struct {
	Id string `json:"id"`
}

type UpdateUserReq struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateUserRes struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
}

type AllUsers struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UsersRes struct {
	Users []AllUsers `json:"users"`
}

type UsersReq struct {
	Limit int32 `json:"limit"`
	Page  int32 `json:"page"`
}
