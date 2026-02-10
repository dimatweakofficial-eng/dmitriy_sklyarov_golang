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
	//проверяем через репозиторий нет ли уже такого емэил в бд
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New("Already exists")
	}
	//шифруем пароль через bcrypt для защиты данных
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	//строим модель для бд, передаем зашифрованный пароль для безопасности
	user := &user.UserModel{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	//обращаемся к репозиторию и создаем запись в бд (пароль зашифрован)
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (servise *AuthServise) Login(email string, password string) (string, error) {
	//проверяем через репозиторий есть ли такой емаил в бд
	existUser, _ := servise.UserRepository.FindByEmail(email)
	if existUser == nil {
		return "", errors.New("Email not exist")
	}
	//проверяем верный ли пароль (шифруем переданный в дто пароль и сравниваем с записью в базе данных)
	err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("Wrong password")
	}
	return existUser.Email, nil
}
