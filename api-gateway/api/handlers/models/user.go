package models

import "api-gateway/api/handlers/models/postModel"

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	UserName  string           `json:"user_name"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	Bio       string           `json:"bio"`
	Website   string           `json:"website"`
	Posts     []postModel.Post `json:"posts"`
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

type UpdateUser struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Bio       string `json:"bio"`
	Website   string `json:"website"`
}

type UpdateUserRes struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Bio       string `json:"bio"`
	Website   string `json:"website"`
	Token     string `json:"token"`
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

type name struct {
}
