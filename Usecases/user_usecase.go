package usecases

import (
	"context"
	"errors"
	repositories "g6/blog-api/Repositories"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Delivery/dto"
	"g6/blog-api/Infrastructure/security"
	"time"
)

type userUsecase struct {
	userRepo   repositories.UserRepository
	ctxtimeout time.Duration
}

func NewUserUsecase(userRepo repositories.UserRepository, timeout time.Duration) *userUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		ctxtimeout: timeout,
	}
}

func (uc *userUsecase) Register(request dto.RegisterRequest) error{
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	if _, err := uc.userRepo.FindByUsernameOrEmail(ctx, request.Username); err == nil {
		return errors.New("username already exists")
	}
	if _, err := uc.userRepo.FindByUsernameOrEmail(ctx, request.Email); err == nil {
		return errors.New("email already exists")
	}
	hashed, _ := security.HashPassword(request.Password)
	now := time.Now()
	user := domain.User{
		Username:       request.Username,	
		Email:          request.Email,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Password:       hashed,
		Role:           domain.RoleUser,
		Bio:            request.Bio,		
		ProfilePicture: request.ProfilePicture,
		CreatedAt:      now,		
		UpdatedAt:      now,
	}
	return uc.userRepo.Create(ctx, user)
}

//Login 
//Logout
func (uc *userUsecase) Logout(userID string) error{
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()
    
	return uc.userRepo.InvalidateTokens(ctx, userID)
}

func (uc *userUsecase) ChangeRole(initiatorRole string, targetUserID string, request dto.ChangeRoleRequest) error{
	ctx, cancel := context.WithTimeout(context.Background(), uc.ctxtimeout)
	defer cancel()

	ir := domain.Role(initiatorRole)
	tr := domain.Role(request.Role)
  
	if ir == domain.RoleAdmin && tr == domain.RoleAdmin{
		return errors.New("only superadmin can promote/ demote admin")
	}
	// Only superadmin and admin can change roles
	if ir != domain.RoleSuperAdmin && ir != domain.RoleAdmin {
		return errors.New("insufficient privileges")
	}
	return uc.userRepo.ChangeRole(ctx, targetUserID, tr)
}

