package models

import "api-gateway/api/handlers/models/postModel"

type UpdateRolReq struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type UpdateRolRes struct {
	Status bool `json:"status"`
}

type AdminUser struct {
	Id        string           `json:"id"`
	UserName  string           `json:"user_name"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	Bio       string           `json:"bio"`
	Website   string           `json:"website"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Posts     []postModel.Post `json:"posts"`
}

type AdminUsersRes struct {
	Users []AdminUser `json:"users"`
}
