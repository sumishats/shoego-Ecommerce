package models


type AdminLogin struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"min=8,max=20"`
}

type AdminDetailsResponse struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"  `
	Lastname  string `json:"lastname" `
	Email     string `json:"email" `
}


//one user row response for admin
type AdminUserResponse struct { 
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Blocked   bool   `json:"blocked"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"created_at"`
}

//all users list  and  pagination ,limit, count..
type AdminUserListResponse struct { 
	Users      []AdminUserResponse `json:"users"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalCount int64               `json:"total_count"`
	TotalPages int                 `json:"total_pages"`
}