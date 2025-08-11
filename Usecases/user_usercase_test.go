package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"g6/blog-api/Infrastructure/security"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// UserUsecaseSuite defines the test suite for UserUsecase
type UserUsecaseSuite struct {
	suite.Suite
	mockUserRepo *domain_mocks.MockIUserRepository
	mockStorage  *domain_mocks.MockStorageService
	usecase      *UserUsecase
	timeout      time.Duration
}

func (s *UserUsecaseSuite) SetupTest() {
	s.mockUserRepo = domain_mocks.NewMockIUserRepository(s.T())
	s.mockStorage = domain_mocks.NewMockStorageService(s.T())
	s.timeout = 5 * time.Second
	s.usecase = &UserUsecase{
		userRepo:       s.mockUserRepo,
		storageService: s.mockStorage,
		ctxtimeout:     s.timeout,
	}
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

func (s *UserUsecaseSuite) TestRegister() {
	s.Run("Success", func() {
		user := &domain.User{
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Username).Return(domain.User{}, errors.New("not found"))
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Email).Return(domain.User{}, errors.New("not found"))
		// hashed password
		hashedPassword, _ := security.HashPassword(user.Password)
		user.Password = hashedPassword
		s.mockUserRepo.On("CreateUser", mock.Anything, user).Return(nil)

		err := s.usecase.Register(user)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("UsernameExists", func() {
		user := &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		existedUser := domain.User{
			Username:   user.Username,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       domain.RoleUser,
			IsVerified: false,
		}
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Username).Return(existedUser, nil)

		err := s.usecase.Register(user)

		s.Error(err)
		s.Equal("username already exists", err.Error())
		s.resetMocks()
	})

	s.Run("EmailExists", func() {
		user := &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		existedUser := domain.User{
			Username:   "abebe_dev",
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       domain.RoleUser,
			IsVerified: false,
		}
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Username).Return(domain.User{}, errors.New("not found"))
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Email).Return(existedUser, nil)

		err := s.usecase.Register(user)

		s.Error(err)
		s.Equal("email already exists", err.Error())
		s.resetMocks()
	})

	s.Run("CreateUserError", func() {
		user := &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Username).Return(domain.User{}, errors.New("not found"))
		s.mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, user.Email).Return(domain.User{}, errors.New("not found"))
		s.mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("create failed"))

		err := s.usecase.Register(user)

		s.Error(err)
		s.Equal("create failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestLogout() {
	s.Run("Success", func() {
		userID := "1"
		s.mockUserRepo.On("InvalidateTokens", mock.Anything, userID).Return(nil)

		err := s.usecase.Logout(userID)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("InvalidateTokensError", func() {
		userID := "1"
		s.mockUserRepo.On("InvalidateTokens", mock.Anything, userID).Return(errors.New("invalidate failed"))

		err := s.usecase.Logout(userID)

		s.Error(err)
		s.Equal("invalidate failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestChangeRole() {
	s.Run("SuccessSuperAdminChangesToAdmin", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleUser, Username: "testuser"}
		newRole := domain.RoleAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)
		s.mockUserRepo.On("ChangeRole", mock.Anything, userID, string(newRole), targetUser.Username).Return(nil)

		err := s.usecase.ChangeRole(string(domain.RoleSuperAdmin), userID, domain.User{Role: newRole})

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("UnauthorizedNonAdmin", func() {
		userID := "1"
		newRole := domain.RoleAdmin

		err := s.usecase.ChangeRole(string(domain.RoleUser), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("unauthorized: only superadmin or admin can change roles", err.Error())
		s.resetMocks()
	})

	s.Run("TargetUserNotFound", func() {
		userID := "1"
		newRole := domain.RoleAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(nil, errors.New("not found"))

		err := s.usecase.ChangeRole(string(domain.RoleSuperAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("target user not found", err.Error())
		s.resetMocks()
	})

	s.Run("SuperAdminRoleByNonSuperAdmin", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleUser}
		newRole := domain.RoleSuperAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)

		err := s.usecase.ChangeRole(string(domain.RoleAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("unauthorized: only superadmin can assign superadmin role", err.Error())
		s.resetMocks()
	})

	s.Run("NoChangeInRole", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleAdmin}
		newRole := domain.RoleAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)

		err := s.usecase.ChangeRole(string(domain.RoleSuperAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("no change in role", err.Error())
		s.resetMocks()
	})

	s.Run("AdminCannotModifySuperAdmin", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleSuperAdmin}
		newRole := domain.RoleAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)

		err := s.usecase.ChangeRole(string(domain.RoleAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("only superadmin can modify superadmin role", err.Error())
		s.resetMocks()
	})

	s.Run("AdminCannotModifyOtherAdmin", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleAdmin}
		newRole := domain.RoleUser
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)

		err := s.usecase.ChangeRole(string(domain.RoleAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("unauthorized: admin cannot modify other admins", err.Error())
		s.resetMocks()
	})

	s.Run("AdminInvalidNewRole", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleUser}
		newRole := domain.RoleSuperAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)

		err := s.usecase.ChangeRole(string(domain.RoleAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("unauthorized: only superadmin can assign superadmin role", err.Error())
		s.resetMocks()
	})

	s.Run("ChangeRoleError", func() {
		userID := "1"
		targetUser := &domain.User{ID: userID, Role: domain.RoleUser, Username: "testuser"}
		newRole := domain.RoleAdmin
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(targetUser, nil)
		s.mockUserRepo.On("ChangeRole", mock.Anything, userID, string(newRole), targetUser.Username).Return(errors.New("change role failed"))

		err := s.usecase.ChangeRole(string(domain.RoleSuperAdmin), userID, domain.User{Role: newRole})

		s.Error(err)
		s.Equal("change role failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestFindByUsernameOrEmail() {
	s.Run("SuccessWithEmail", func() {
		identifier := "test@example.com"
		user := &domain.User{ID: "1", Email: identifier}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, identifier).Return(user, nil)

		result, err := s.usecase.FindByUsernameOrEmail(context.Background(), identifier)

		s.NoError(err)
		s.Equal(user, result)
		s.resetMocks()
	})

	s.Run("SuccessWithUsername", func() {
		identifier := "testuser"
		user := &domain.User{ID: "1", Username: identifier}
		s.mockUserRepo.On("GetUserByUsername", mock.Anything, identifier).Return(user, nil)

		result, err := s.usecase.FindByUsernameOrEmail(context.Background(), identifier)

		s.NoError(err)
		s.Equal(user, result)
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		identifier := "test@example.com"
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, identifier).Return(nil, errors.New("not found"))

		result, err := s.usecase.FindByUsernameOrEmail(context.Background(), identifier)

		s.Error(err)
		s.Nil(result)
		s.Equal("not found", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestFindUserByID() {
	s.Run("Success", func() {
		userID := "1"
		user := &domain.User{ID: userID}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)

		result, err := s.usecase.FindUserByID(userID)

		s.NoError(err)
		s.Equal(user, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		userID := "1"
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(nil, errors.New("not found"))

		result, err := s.usecase.FindUserByID(userID)

		s.Error(err)
		s.Nil(result)
		s.Equal("not found", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestGetUserByEmail() {
	s.Run("Success", func() {
		email := "test@example.com"
		user := &domain.User{Email: email}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

		result, err := s.usecase.GetUserByEmail(email)

		s.NoError(err)
		s.Equal(user, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		email := "test@example.com"
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("not found"))

		result, err := s.usecase.GetUserByEmail(email)

		s.Error(err)
		s.Nil(result)
		s.Equal("not found", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestUpdateUser() {
	s.Run("Success", func() {
		userID := "1"
		user := &domain.User{ID: userID, Username: "testuser"}
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
			return u.ID == userID && u.Username == "testuser" && !u.UpdatedAt.IsZero()
		})).Return(nil)

		result, err := s.usecase.UpdateUser(userID, user)

		s.NoError(err)
		s.Equal(user, result)
		s.False(result.UpdatedAt.IsZero())
		s.resetMocks()
	})

	s.Run("UpdateError", func() {
		userID := "1"
		user := &domain.User{ID: userID}
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.Anything).Return(errors.New("update failed"))

		result, err := s.usecase.UpdateUser(userID, user)

		s.Error(err)
		s.Nil(result)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestUpdateProfile() {
	s.Run("SuccessWithAvatar", func() {
		userID := "1"
		user := &domain.User{ID: userID, Username: "testuser", Email: "test@example.com"}
		update := domain.UserProfileUpdate{
			Bio:        "Updated bio",
			FirstName:  "Test",
			LastName:   "User",
			AvatarData: []byte("fake-image-data"),
		}
		fileName := "avatar.jpg"
		avatarURL := "http://storage.com/avatar.jpg"
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		s.mockStorage.On("UploadFile", mock.Anything, fileName, update.AvatarData).Return(avatarURL, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
			return u.ID == userID &&
				u.Bio == update.Bio &&
				u.FirstName == update.FirstName &&
				u.LastName == update.LastName &&
				u.AvatarURL == avatarURL
		})).Return(nil)

		result, err := s.usecase.UpdateProfile(userID, update, fileName)

		s.NoError(err)
		s.Equal(userID, result.ID)
		s.Equal(update.Bio, result.Bio)
		s.Equal(update.FirstName, result.FirstName)
		s.Equal(update.LastName, result.LastName)
		s.Equal(avatarURL, result.AvatarURL)
		s.resetMocks()
	})

	s.Run("SuccessWithoutAvatar", func() {
		userID := "1"
		user := &domain.User{ID: userID, Username: "testuser", Email: "test@example.com"}
		update := domain.UserProfileUpdate{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
			return u.ID == userID &&
				u.Bio == update.Bio &&
				u.FirstName == update.FirstName &&
				u.LastName == update.LastName &&
				u.AvatarURL == ""
		})).Return(nil)

		result, err := s.usecase.UpdateProfile(userID, update, "")

		s.NoError(err)
		s.Equal(userID, result.ID)
		s.Equal(update.Bio, result.Bio)
		s.Equal(update.FirstName, result.FirstName)
		s.Equal(update.LastName, result.LastName)
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		userID := "1"
		update := domain.UserProfileUpdate{Bio: "Updated bio"}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)

		result, err := s.usecase.UpdateProfile(userID, update, "")

		s.Error(err)
		s.Nil(result)
		s.Equal(domain.ErrUserNotFound, err)
		s.resetMocks()
	})

	s.Run("UploadFileError", func() {
		userID := "1"
		user := &domain.User{ID: userID, Username: "testuser"}
		update := domain.UserProfileUpdate{
			Bio:        "Updated bio",
			AvatarData: []byte("fake-image-data"),
		}
		fileName := "avatar.jpg"
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		s.mockStorage.On("UploadFile", mock.Anything, fileName, update.AvatarData).Return("", errors.New("upload failed"))

		result, err := s.usecase.UpdateProfile(userID, update, fileName)

		s.Error(err)
		s.Nil(result)
		s.Contains(err.Error(), "failed to upload avatar")
		s.resetMocks()
	})

	s.Run("UpdateUserError", func() {
		userID := "1"
		user := &domain.User{ID: userID, Username: "testuser"}
		update := domain.UserProfileUpdate{Bio: "Updated bio"}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.Anything).Return(errors.New("update failed"))

		result, err := s.usecase.UpdateProfile(userID, update, "")

		s.Error(err)
		s.Nil(result)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) TestChangePassword() {
	s.Run("Success", func() {
		userID := "1"
		oldPassword := "oldpassword123"
		newPassword := "newpassword123"

		hashedOldPassword, _ := security.HashPassword(oldPassword)
		user := &domain.User{ID: userID, Password: hashedOldPassword}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, user).Return(nil)

		err := s.usecase.ChangePassword(userID, oldPassword, newPassword)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		userID := "1"
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(nil, errors.New("not found"))

		err := s.usecase.ChangePassword(userID, "oldpassword123", "newpassword123")

		s.Error(err)
		s.Equal("not found", err.Error())
		s.resetMocks()
	})

	s.Run("InvalidOldPassword", func() {
		userID := "1"
		hashedPassword, _ := security.HashPassword("correctpassword")
		user := &domain.User{ID: userID, Password: hashedPassword}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)

		err := s.usecase.ChangePassword(userID, "wrongpassword", "newpassword123")

		s.Error(err)
		s.Equal("invalid old password", err.Error())
		s.resetMocks()
	})

	s.Run("HashPasswordError", func() {
		userID := "1"
		oldPassword := "oldpassword123"
		hashedPassword, _ := security.HashPassword(oldPassword)
		user := &domain.User{ID: userID, Password: hashedPassword}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)
		hashedPassword, _ = security.HashPassword(oldPassword)
		err := security.ValidatePassword(hashedPassword, oldPassword)
		s.Require().NoError(err)
		// mock hash password cause error
		newPassword := strings.Repeat("z", 1000) // Simulate a long password that might cause hashing to fail
		err = s.usecase.ChangePassword(userID, oldPassword, newPassword)

		s.Error(err)
		s.resetMocks()
	})

	s.Run("UpdateUserError", func() {
		userID := "1"
		oldPassword := "oldpassword123"
		newPassword := "newpassword123"
		hashedPassword, _ := security.HashPassword(oldPassword)
		user := &domain.User{ID: userID, Password: hashedPassword}
		s.mockUserRepo.On("FindUserByID", mock.Anything, userID).Return(user, nil)

		s.mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.Anything).Return(errors.New("update failed"))

		err := s.usecase.ChangePassword(userID, oldPassword, newPassword)

		s.Error(err)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})
}

func (s *UserUsecaseSuite) resetMocks() {
	s.mockUserRepo.ExpectedCalls = nil
	s.mockUserRepo.Calls = nil
	s.mockStorage.ExpectedCalls = nil
	s.mockStorage.Calls = nil
}
