package route

type CreateUserGroupDto struct {
	GroupName string `json:"name"`
	GroupDesc string `json:"description"`
	GroupType string `json:"type"`
}

type UpdateUserRoleDto struct {
	Role string `json:"role"`
}
