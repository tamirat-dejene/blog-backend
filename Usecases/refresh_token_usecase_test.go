package usecases

import (
	"errors"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// RefreshTokenUsecaseSuite defines the test suite for RefreshTokenUsecase
type RefreshTokenUsecaseSuite struct {
	suite.Suite
	mockRepo *domain_mocks.MockIRefreshTokenRepository
	usecase  *RefreshTokenUsecase
}

func (s *RefreshTokenUsecaseSuite) SetupTest() {
	s.mockRepo = domain_mocks.NewMockIRefreshTokenRepository(s.T())
	s.usecase = &RefreshTokenUsecase{
		Repo: s.mockRepo,
	}
}

func TestRefreshTokenUsecaseSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenUsecaseSuite))
}

func (s *RefreshTokenUsecaseSuite) TestFindByToken() {
	s.Run("Success", func() {
		token := "refresh-token"
		expectedToken := &domain.RefreshToken{
			Token:     token,
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("FindByToken", mock.Anything, token).Return(expectedToken, nil)

		result, err := s.usecase.FindByToken(token)

		s.NoError(err)
		s.Equal(expectedToken, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		token := "refresh-token"
		s.mockRepo.On("FindByToken", mock.Anything, token).Return(nil, errors.New("token not found"))

		result, err := s.usecase.FindByToken(token)

		s.Error(err)
		s.Nil(result)
		s.Equal("token not found", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		token := "refresh-token"
		s.mockRepo.On("FindByToken", mock.Anything, token).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(nil, errors.New("context deadline exceeded"))

		result, err := s.usecase.FindByToken(token)

		s.Error(err)
		s.Nil(result)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) TestSave() {
	s.Run("Success", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("Save", mock.Anything, token).Return(nil)

		err := s.usecase.Save(token)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("SaveError", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("Save", mock.Anything, token).Return(errors.New("save failed"))

		err := s.usecase.Save(token)

		s.Error(err)
		s.Equal("save failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("Save", mock.Anything, token).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(errors.New("context deadline exceeded"))

		err := s.usecase.Save(token)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) TestDeleteByUserID() {
	s.Run("Success", func() {
		userID := "1"
		s.mockRepo.On("DeleteByUserID", mock.Anything, userID).Return(nil)

		err := s.usecase.DeleteByUserID(userID)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("DeleteError", func() {
		userID := "1"
		s.mockRepo.On("DeleteByUserID", mock.Anything, userID).Return(errors.New("delete failed"))

		err := s.usecase.DeleteByUserID(userID)

		s.Error(err)
		s.Equal("delete failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		userID := "1"
		s.mockRepo.On("DeleteByUserID", mock.Anything, userID).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(errors.New("context deadline exceeded"))

		err := s.usecase.DeleteByUserID(userID)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) TestReplaceToken() {
	s.Run("Success", func() {
		token := &domain.RefreshToken{
			Token:     "new-refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("ReplaceTokenByUserID", mock.Anything, token).Return(nil)

		err := s.usecase.ReplaceToken(token)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("ReplaceError", func() {
		token := &domain.RefreshToken{
			Token:     "new-refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("ReplaceTokenByUserID", mock.Anything, token).Return(errors.New("replace failed"))

		err := s.usecase.ReplaceToken(token)

		s.Error(err)
		s.Equal("replace failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		token := &domain.RefreshToken{
			Token:     "new-refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("ReplaceTokenByUserID", mock.Anything, token).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(errors.New("context deadline exceeded"))

		err := s.usecase.ReplaceToken(token)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) TestRevokedToken() {
	s.Run("Success", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("RevokeToken", mock.Anything, token.Token).Return(nil)

		err := s.usecase.RevokedToken(token)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("RevokeError", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("RevokeToken", mock.Anything, token.Token).Return(errors.New("revoke failed"))

		err := s.usecase.RevokedToken(token)

		s.Error(err)
		s.Equal("revoke failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		token := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("RevokeToken", mock.Anything, token.Token).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(errors.New("context deadline exceeded"))

		err := s.usecase.RevokedToken(token)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) TestFindByUserID() {
	s.Run("Success", func() {
		userID := "1"
		expectedToken := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    userID,
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRepo.On("FindTokenByUserID", mock.Anything, userID).Return(expectedToken, nil)

		result, err := s.usecase.FindByUserID(userID)

		s.NoError(err)
		s.Equal(expectedToken, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		userID := "1"
		s.mockRepo.On("FindTokenByUserID", mock.Anything, userID).Return(nil, errors.New("token not found"))

		result, err := s.usecase.FindByUserID(userID)

		s.Error(err)
		s.Nil(result)
		s.Equal("token not found", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		userID := "1"
		s.mockRepo.On("FindTokenByUserID", mock.Anything, userID).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(nil, errors.New("context deadline exceeded"))

		result, err := s.usecase.FindByUserID(userID)

		s.Error(err)
		s.Nil(result)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *RefreshTokenUsecaseSuite) resetMocks() {
	s.mockRepo.ExpectedCalls = nil
	s.mockRepo.Calls = nil
}
