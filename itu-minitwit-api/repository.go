package main

import (
	"gorm.io/gorm"
)

//GetUserByUsername ...
func GetUserByUsername(username string, db *gorm.DB) (User, error) {
	user := User{}
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}

//GetUserById ...
func GetUserById(id uint, db *gorm.DB) (User, error) {
	user := User{}
	result := db.First(&user, id)
	return user, result.Error
}
