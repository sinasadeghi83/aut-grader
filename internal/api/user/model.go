package user

import "github.com/sinasadeghi83/aut-grader/internal/api/platform/model"

type User struct {
	model.Model
	Username string `json:"username" gorm:"unique;type:varchar(25)"`
	Password string `json:"-" gorm:",type:varchar(25)"`
}

func (User) TableName() string {
	return "users"
}
