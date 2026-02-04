package auth

import (
	"demo-1/intern/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthServise struct {
	UserRepository *user.UserRepositiry
}

func NewAuthService(repo *user.UserRepositiry) *AuthServise {
	return &AuthServise{
		UserRepository: repo,
	}
}

func (service *AuthServise) Register(email string, password string, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New("Already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.UserModel{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (servise *AuthServise) Login(email string, password string) (string, error) {
	existUser, _ := servise.UserRepository.FindByEmail(email)
	if existUser == nil {
		return "", errors.New("Email not exist")
	}
	err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("Wrong password")
	}
	return string(email), nil
}
