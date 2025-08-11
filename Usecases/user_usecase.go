package usecases

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"regexp"
	"time"
)

type UserUsecase struct {
	userRepo       domain.IUserRepository
	storageService domain.StorageService
	ctxtimeout     time.Duration
}

func NewUserUsecase(userRepo domain.IUserRepository, storageService domain.StorageService, timeout time.Duration) domain.IUserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		storageService: storageService,
		ctxtimeout:     timeout,
	}
}

func (uc *UserUsecase) Register(request *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	request.Role = domain.RoleUser // Default role is User
	user, err := uc.userRepo.FindByUsernameOrEmail(ctx, request.Username)
	if err == nil && (user != domain.User{}) {
		return errors.New("username already exists")
	}

	user, err = uc.userRepo.FindByUsernameOrEmail(ctx, request.Email)
	if err == nil && (user != domain.User{}) {
		return errors.New("email already exists")
	}
	hashed, _ := security.HashPassword(request.Password)
	request.Password = hashed
	request.IsVerified = false
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	return uc.userRepo.CreateUser(ctx, request)
}

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

func (uc *UserUsecase) UpdateProfile(userID string, update domain.UserProfileUpdate, fileName string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user
	user, err := uc.userRepo.FindUserByID(ctx, userID)
	if err == domain.ErrNotFound {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(update.FirstName, update.LastName)
	// Handle avatar upload
	if len(update.AvatarData) > 0 {
		avatarURL, err := uc.storageService.UploadFile(ctx, fileName, update.AvatarData)
		if err != nil {
			return nil, fmt.Errorf("failed to upload avatar: %w", err)
		}

		// Wait for the upload to complete
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue if no timeout or cancellation
		}

		user.AvatarURL = avatarURL
	}
	fmt.Println(user)

	// Apply updates
	if update.Bio != "" {
		user.Bio = update.Bio
	}
	if update.FirstName != "" {
		user.FirstName = update.FirstName
	}
	if update.LastName != "" {
		user.LastName = update.LastName
	}

	// Update in repository
	if err := uc.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) ChangePassword(userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	// Find user
	user, err := uc.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check old password
	if err := security.ValidatePassword(user.Password, oldPassword); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	return uc.userRepo.UpdateUser(ctx, user.ID, user)
}
