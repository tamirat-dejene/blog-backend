package dto

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	FirstName	 string `json:"first_name"`
	LastName      string `json:"last_name"`
	Password	   string `json:"password" binding:"required,min=8"`
	Bio           string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
}

//LoginRequest
//AuthResponse

type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}