package domain
import "time"

type Role string

const(
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleSuperAdmin Role = "super_admin"
)

type TokenPair struct{
	AccessToken string
	RefreshToken string
}

type User struct {
	ID 	  string	
	Username string
	Email    string	
	FirstName string
	LastName  string
	Password string
	Role Role
	Tokens TokenPair
	Bio string
	ProfilePicture string
	CreatedAt time.Time	
	UpdatedAt time.Time
}
