package models

type ListRolesResponse struct {
	Roles []string `json:"roles"`
	Obj   []string `json:"obj"`
	Act   []string `json:"act"`
}
type CreateUserRoleRequest struct {
	Userid string `json:"userid"`
	Role   string `json:"role"`
}

type AddRole struct {
	Role   string `json:"role"`
	Url    string `json:"url"`
	Method string `json:"method"`
}
