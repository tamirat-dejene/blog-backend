package usecases

import (
	"context"
	"g6/blog-api/Domain"
	"time"
)

type UserUsecase struct {
	UserRepository Domain.IUserRepository
}

func NewUserUsecase(userRepository Domain.IUserRepository) Domain.IUserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (usecase *UserUsecase) GetAllUsers() ([]*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	users, err := usecase.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (usecase *UserUsecase) CreateUser(user *Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := usecase.UserRepository.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (usecase *UserUsecase) GetUserByUsername(username string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := usecase.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetUserByEmail(email string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := usecase.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
