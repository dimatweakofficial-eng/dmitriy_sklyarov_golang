package user

import "demo-1/pkg/db"

type UserRepositiry struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepositiry {
	return &UserRepositiry{
		Database: db,
	}
}

func (repo *UserRepositiry) Create(newUser *UserModel) (*UserModel, error) {
	result := repo.Database.DB.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}

func (repo *UserRepositiry) FindByEmail(email string) (*UserModel, error) {
	var user UserModel
	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
