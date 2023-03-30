package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sdu.store/server"
)

type User struct {
	gorm.Model
	Email    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
	Is_admin bool   `json:"is_admin"`
	Is_staff bool   `json:"is_staff"`
}

var Users []User

func GetUserByUsername(username string) (*User, error) {
	user := User{}
	if err := server.DB.Where("username", username).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := User{}

	if err := server.DB.Where("id", id).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *User) IsAdmin() bool {
	return user.Is_admin
}

func (user *User) IsStaff() bool {
	return user.Is_staff
}

func (user *User) Delete() error {
	if err := server.DB.Where("ID=?", user.ID).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) Update() error {
	isStaff := user.Is_staff
	isAdmin := user.Is_admin
	if err := server.DB.First(user).Error; err != nil {
		return err
	}
	user.Is_admin = isAdmin
	user.Is_staff = isStaff || isAdmin
	if err := server.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (this *User) BeforeCreate(tx *gorm.DB) (err error) {
	this.Password, err = hashPassword(this.Password)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
