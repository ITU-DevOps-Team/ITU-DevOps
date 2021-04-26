package main

import (
	"fmt"
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

	db.Where(
		"flagged = ?", 0).Limit(numberOfPosts).Order(
		"pub_date desc").Find(&posts)

	return ConvertMessagesToViewPosts(posts, db)
}

func GetPostsByUser(username string, db *gorm.DB) []ViewPost{
	var posts []Message
	db.Table(
		"messages").Order(
		"messages.pub_date desc").Select(
		"users.username, messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged").Joins(
		"join users on users.user_id = messages.author_id").Where(
		"messages.flagged = 0 AND users.username = ?", username).Scan(&posts)

	return ConvertMessagesToViewPosts(posts, db)
}

func CheckIfUserIsFollowed(who string, whom string, db *gorm.DB) bool {
	if whom == "" {
		return false
	}
	if who == "" {
		return false
	}

	whomUser, err := GetUserByUsername(whom, db)
	checkErr(err)
	whoUser, err := GetUserByUsername(who, db)
	checkErr(err)

	follower := []Follower{}
	result := db.Table(
		"followers").Where(
			"who_id = ? AND whom_id = ?", whoUser.UserID, whomUser.UserID).Scan(&follower)
	fmt.Println(result)
	return len(follower) > 1
}

func ConvertMessagesToViewPosts(messages []Message, db *gorm.DB) []ViewPost{
	var postSlice []ViewPost

	for _, message := range messages {
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