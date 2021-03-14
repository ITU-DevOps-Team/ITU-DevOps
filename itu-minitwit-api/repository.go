package main

import (
	"gorm.io/gorm"
)

func GetLatest(db *gorm.DB) (Latest, error) {
	latest := Latest{}
	result := db.Last(&latest)
	return latest, result.Error
}

func AddLatest(latest Latest, db *gorm.DB) error {
	result := db.Create(&latest)
	return result.Error
}

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
