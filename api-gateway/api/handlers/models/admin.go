package models

type UpdateRolReq struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type UpdateRolRes struct {
	Status bool `json:"status"`
}
