package dto

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=6,max=100"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user superadmin"`
}