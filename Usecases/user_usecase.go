package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"regexp"
	"time"
)

type userUsecase struct {
	userRepo   domain.IUserRepository
	ctxtimeout time.Duration
}

func NewUserUsecase(userRepo domain.IUserRepository, timeout time.Duration) domain.IUserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		ctxtimeout: timeout,
	}
}

func (uc *userUsecase) Register(request *domain.User) error {
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
func (uc *userUsecase) Logout(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	return uc.userRepo.InvalidateTokens(ctx, userID)
}

func (uc *userUsecase) ChangeRole(initiatorRole string, targetUserID string, request domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	if domain.UserRole(initiatorRole) == domain.RoleAdmin && domain.UserRole(request.Role) == domain.RoleAdmin {
		return errors.New("only superadmin can promote/ demote admin")
	}
	// Only superadmin and admin can change roles
	if domain.UserRole(initiatorRole) != domain.RoleSuperAdmin && domain.UserRole(initiatorRole) != domain.RoleAdmin {
		return errors.New("insufficient privileges")
	}
	return uc.userRepo.ChangeRole(ctx, targetUserID, string(domain.UserRole(request.Role)), request.Username)
}

// find user by username or id
func (uc *userUsecase) FindByUsernameOrEmail(ctx context.Context, identifier string) (*domain.User, error) {
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

func (uc *userUsecase) FindUserByID(uid string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()
	return uc.userRepo.FindUserByID(ctx, uid)
}

func (uc *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()
	return uc.userRepo.GetUserByEmail(ctx, email)
}

func (uc *userUsecase) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	user.UpdatedAt = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, id, user); err != nil {
		return nil, err
	}
	return user, nil
}
