package main

import (
	"gorm.io/gorm"
)

func GetUserByUsername(username string, db *gorm.DB) (User, error) {
	user := User{}
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}

func GetUserById(id uint, db *gorm.DB) (User, error) {
	user := User{}
	result := db.First(&user, id)
	return user, result.Error
}
