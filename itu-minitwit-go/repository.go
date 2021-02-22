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

func GetPublicPosts(numberOfPosts int, db *gorm.DB) []ViewPost {
	var posts []Message

	db.Where("flagged = ?", 0).Limit(numberOfPosts).Order("pub_date desc").Find(&posts)

	var postSlice []ViewPost

	for _, message := range posts {
		user, _ := GetUserById(message.Author_id, db)

		post := ViewPost{
			Username:      user.Username,
			Message_id:	   message.Message_id,
			Author_id:	   message.Author_id,
			Text:          message.Text,
			Pub_date:	   getTimeFromTimestamp(message.Pub_date),
			Flagged:       message.Flagged,
		}
		postSlice = append(postSlice, post)
	}
	return postSlice
}
