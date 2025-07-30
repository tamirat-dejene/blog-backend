package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type UserUsecase struct {
	UserRepository domain.IUserRepository
}

func NewUserUsecase(userRepository domain.IUserRepository) domain.IUserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (usecase *UserUsecase) GetAllUsers() ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	users, err := usecase.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (usecase *UserUsecase) CreateUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := usecase.UserRepository.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (usecase *UserUsecase) GetUserByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := usecase.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := usecase.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
