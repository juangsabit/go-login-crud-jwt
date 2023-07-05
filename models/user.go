package models

import (
	"errors"
	"go-postgres/utils/token"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	FullName  string `gorm:"type:varchar(255)" json:"fullname"`
	Username  string `gorm:"type:varchar(255)" binding:"required" json:"username"`
	Password  string `gorm:"type:varchar(255)" binding:"required" json:"password"`
	Email     string `gorm:"type:varchar(255)" json:"email"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Activitys []Activity `gorm:"foreignKey:UserID;references:ID"`
}

type Activity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	UserID      uint   `json:"user_id"`
	Information string `gorm:"type:varchar(255)" json:"information"`
	TableName   string `gorm:"type:varchar(50)" json:"table_name"`
	TableID     uint   `json:"table_id"`
	CreatedAt   *time.Time
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func LoginCheck(username string, password string) (string, error) {

	var err error
	u := User{}
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	match := VerifyPassword(password, u.Password)
	if !match {
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)); err != nil {
			return "", err
		}
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

type UserLogin struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" binding:"required" json:"username"`
	Password string `gorm:"size:255;not null;" binding:"required" json:"password"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New(strconv.Itoa(int(uid)) + "User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}
