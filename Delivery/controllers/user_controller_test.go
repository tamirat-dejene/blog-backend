package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// UserControllerSuite defines the test suite for UserController
type UserControllerSuite struct {
	suite.Suite
	mockUserUsecase *domain_mocks.MockIUserUsecase
	handler         *UserController
	validate        *validator.Validate
}

// SetupTest initializes the mocks and handler before each test
func (s *UserControllerSuite) SetupTest() {
	s.mockUserUsecase = domain_mocks.NewMockIUserUsecase(s.T())
	s.handler = &UserController{
		uc: s.mockUserUsecase,
	}
	s.validate = validator.New()
}

// TestUserControllerSuite runs the test suite
func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}

// TestUpdateProfile tests the UpdateProfile method
func (s *UserControllerSuite) TestUpdateProfile() {
	s.Run("SuccessWithAvatar", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}
		updatedUser := &domain.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			FirstName: updateRequest.FirstName,
			LastName:  updateRequest.LastName,
			Bio:       updateRequest.Bio,
			AvatarURL: "avatar.jpg",
		}
		avatarData := []byte("fake-image-data")
		fileName := "avatar.jpg"

		s.mockUserUsecase.On("UpdateProfile", userID, mock.MatchedBy(func(u domain.UserProfileUpdate) bool {
			return u.Bio == updateRequest.Bio &&
				u.FirstName == updateRequest.FirstName &&
				u.LastName == updateRequest.LastName &&
				bytes.Equal(u.AvatarData, avatarData)
		}), fileName).Return(updatedUser, nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		part, _ := writer.CreateFormFile("avatar_file", fileName)
		_, _ = part.Write(avatarData)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(updatedUser.Username, response.Username)
		s.Equal(updatedUser.Email, response.Email)
		s.Equal(updatedUser.FirstName, response.FirstName)
		s.Equal(updatedUser.LastName, response.LastName)
		s.Equal(updatedUser.Bio, response.Bio)
		s.Equal(updatedUser.AvatarURL, response.AvatarURL)
		s.resetMocks()
	})

	s.Run("SuccessWithoutAvatar", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}
		updatedUser := &domain.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			FirstName: updateRequest.FirstName,
			LastName:  updateRequest.LastName,
			Bio:       updateRequest.Bio,
		}

		s.mockUserUsecase.On("UpdateProfile", userID, mock.MatchedBy(func(u domain.UserProfileUpdate) bool {
			return u.Bio == updateRequest.Bio &&
				u.FirstName == updateRequest.FirstName &&
				u.LastName == updateRequest.LastName &&
				u.AvatarData == nil
		}), "").Return(updatedUser, nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(updatedUser.Username, response.Username)
		s.Equal(updatedUser.Email, response.Email)
		s.Equal(updatedUser.FirstName, response.FirstName)
		s.Equal(updatedUser.LastName, response.LastName)
		s.Equal(updatedUser.Bio, response.Bio)
		s.resetMocks()
	})

	s.Run("Unauthorized", func() {
		req := httptest.NewRequest(http.MethodPost, "/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(domain.ErrUnauthorized.Error(), response["error"])
		s.resetMocks()
	})

	s.Run("InvalidFormData", func() {
		req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader("invalid form data"))
		req.Header.Set("Content-Type", "multipart/form-data")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "1")

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid form data", response["error"])
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       strings.Repeat("a", 1001), // Too long
			FirstName: "",                        // Missing
			LastName:  "",                        // Missing
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Field validation")
		s.resetMocks()
	})

	s.Run("FailedToOpenFile", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		_, _ = writer.CreateFormFile("avatar_file", "avatar.jpg") // No data written
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		// Mock FormFile to return an invalid file
		c.Request.MultipartForm = &multipart.Form{
			File: map[string][]*multipart.FileHeader{
				"avatar_file": {{Filename: "avatar.jpg"}},
			},
		}

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Error uploading file")
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}

		s.mockUserUsecase.On("UpdateProfile", userID, mock.MatchedBy(func(u domain.UserProfileUpdate) bool {
			return u.Bio == updateRequest.Bio &&
				u.FirstName == updateRequest.FirstName &&
				u.LastName == updateRequest.LastName &&
				u.AvatarData == nil
		}), "").Return(nil, domain.ErrUserNotFound)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("User not found", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidFileFormat", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}
		avatarData := []byte("fake-image-data")
		fileName := "avatar.invalid"

		s.mockUserUsecase.On("UpdateProfile", userID, mock.MatchedBy(func(u domain.UserProfileUpdate) bool {
			return u.Bio == updateRequest.Bio &&
				u.FirstName == updateRequest.FirstName &&
				u.LastName == updateRequest.LastName &&
				bytes.Equal(u.AvatarData, avatarData)
		}), fileName).Return(nil, domain.ErrInvalidFile)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		part, _ := writer.CreateFormFile("avatar_file", fileName)
		_, _ = part.Write(avatarData)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid file format", response["error"])
		s.resetMocks()
	})

	s.Run("UpdateError", func() {
		userID := "1"
		updateRequest := dto.UserUpdateProfileRequest{
			Bio:       "Updated bio",
			FirstName: "Test",
			LastName:  "User",
		}

		s.mockUserUsecase.On("UpdateProfile", userID, mock.MatchedBy(func(u domain.UserProfileUpdate) bool {
			return u.Bio == updateRequest.Bio &&
				u.FirstName == updateRequest.FirstName &&
				u.LastName == updateRequest.LastName &&
				u.AvatarData == nil
		}), "").Return(nil, errors.New("update failed"))

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("bio", updateRequest.Bio)
		_ = writer.WriteField("first_name", updateRequest.FirstName)
		_ = writer.WriteField("last_name", updateRequest.LastName)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/profile", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.UpdateProfile(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("update failed", response["error"])
		s.resetMocks()
	})
}

// TestChangePassword tests the ChangePassword method
func (s *UserControllerSuite) TestChangePassword() {
	s.Run("Success", func() {
		userID := "1"
		changePasswordRequest := dto.ChangePasswordRequest{
			OldPassword: "oldpassword123",
			NewPassword: "newpassword123",
		}

		s.mockUserUsecase.On("ChangePassword", userID, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword).Return(nil)

		body, _ := json.Marshal(changePasswordRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.ChangePassword(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Password changed successfully", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "1")

		s.handler.ChangePassword(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid request", response["error"])
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		userID := "1"
		changePasswordRequest := dto.ChangePasswordRequest{
			OldPassword: "",   // Missing
			NewPassword: "sh", // Too short
		}

		body, _ := json.Marshal(changePasswordRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.ChangePassword(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
		s.resetMocks()
	})

	s.Run("Unauthorized", func() {
		changePasswordRequest := dto.ChangePasswordRequest{
			OldPassword: "oldpassword123",
			NewPassword: "newpassword123",
		}

		body, _ := json.Marshal(changePasswordRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.ChangePassword(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Unauthorized", response["error"])
		s.resetMocks()
	})

	s.Run("ChangePasswordError", func() {
		userID := "1"
		changePasswordRequest := dto.ChangePasswordRequest{
			OldPassword: "oldpassword123",
			NewPassword: "newpassword123",
		}

		s.mockUserUsecase.On("ChangePassword", userID, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword).Return(errors.New("invalid old password"))

		body, _ := json.Marshal(changePasswordRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", userID)

		s.handler.ChangePassword(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("invalid old password", response["error"])
		s.resetMocks()
	})
}

func (s *UserControllerSuite) resetMocks() {
	s.mockUserUsecase.ExpectedCalls = nil
	s.mockUserUsecase.Calls = nil
}
