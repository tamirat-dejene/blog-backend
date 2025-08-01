package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"regexp"
	"time"
)

type UserUsecase struct {
	userRepo   domain.IUserRepository
	ctxtimeout time.Duration
}

func NewUserUsecase(userRepo domain.IUserRepository, timeout time.Duration) domain.IUserUsecase {
	return &UserUsecase{
		userRepo:   userRepo,
		ctxtimeout: timeout,
	}
}

func (uc *UserUsecase) Register(request *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	request.Role = domain.RoleUser // Default role is User
	if _, err := uc.userRepo.FindByUsernameOrEmail(ctx, request.Username); err == nil {
		return errors.New("username already exists")
	}
	if _, err := uc.userRepo.FindByUsernameOrEmail(ctx, request.Email); err == nil {
		return errors.New("email already exists")
	}
	hashed, _ := security.HashPassword(request.Password)
	request.Password = hashed
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	return uc.userRepo.CreateUser(ctx, request)
}

// Login
// Logout
func (uc *UserUsecase) Logout(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	return uc.userRepo.InvalidateTokens(ctx, userID)
}

func (uc *UserUsecase) ChangeRole(initiatorRole string, targetUserID string, request domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	initRole := domain.UserRole(initiatorRole)
	newRole := request.Role

	// only super admin or admin can change roles
	if initRole != domain.RoleSuperAdmin && initRole != domain.RoleAdmin {
		return errors.New("unauthorized: only superadmin or admin can change roles")
	}

	target, err := uc.userRepo.FindUserByID(ctx, targetUserID)
	if err != nil {
		return errors.New("target user not found")
	}

	//only super admin can change to super admin
	if newRole == domain.RoleSuperAdmin && initRole != domain.RoleSuperAdmin {
		return errors.New("unauthorized: only superadmin can assign superadmin role")
	}
	if newRole == target.Role {
		return errors.New("no change in role")
	}

	if target.Role == domain.RoleSuperAdmin && initRole != domain.RoleSuperAdmin {
		return errors.New("only superadmin can modify superadmin role")
	}
	if initRole == domain.RoleAdmin {
		switch target.Role {
		case domain.RoleAdmin:
			return errors.New("unauthorized: admin cannot modify other admins")
		case domain.RoleSuperAdmin:
			return errors.New("unauthorized: admin cannot modify superadmin")
		case domain.RoleUser:
			if newRole != domain.RoleAdmin {
				return errors.New("unauthorized: admin can only promote users to admin")
			}
		}
	}

	// Proceed with changing the role
	return uc.userRepo.ChangeRole(ctx, targetUserID, string(newRole), target.Username)
}

// find user by username or id
func (uc *UserUsecase) FindByUsernameOrEmail(ctx context.Context, identifier string) (*domain.User, error) {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, identifier)
	var user *domain.User
	var err error
	if isEmail {
		user, err = uc.userRepo.GetUserByEmail(ctx, identifier)
	} else {
		user, err = uc.userRepo.GetUserByUsername(ctx, identifier)
	}
	return user, err
}

func (uc *UserUsecase) FindUserByID(uid string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()
	return uc.userRepo.FindUserByID(ctx, uid)
}

func (uc *UserUsecase) GetUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()
	return uc.userRepo.GetUserByEmail(ctx, email)
}

func (uc *UserUsecase) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	user.UpdatedAt = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, id, user); err != nil {
		return nil, err
	}
	return user, nil
}
