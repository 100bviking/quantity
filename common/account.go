package common

import (
	"quantity/common/db"
)

type User struct {
	ID        int64  `gorm:"column:id"`
	UserName  string `gorm:"column:user_name"`
	ApiKey    string `gorm:"column:api_key"`
	ApiSecret string `gorm:"column:api_secret"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUser() (user *User, err error) {
	if err != nil {
		return
	}
	user = new(User)
	err = db.AccountDB.Model(user).Where("id = 1").Take(user).Error
	return
}
