package user

import (
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo UserRepo) WithDB(db *gorm.DB) *UserRepo {
	repo.db = db
	return &repo
}

func (repo *UserRepo) Create(newUser User) (User, error) {
	result := repo.db.Create(&newUser)
	return newUser, result.Error
}

func (repo *UserRepo) FindById(id uint) (*User, error) {
	var user User
	result := repo.db.First(&user, id)
	return &user, result.Error
}

func (repo *UserRepo) FindByUsername(username string) (*User, error) {
	var user User
	result := repo.db.Where("username = ?", username).First(&user)
	return &user, result.Error
}
