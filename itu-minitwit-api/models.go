package main

//import ("gorm.io/gorm")

type User struct {
	UserID   uint `gorm:"primaryKey;autoIncrement"`
	Username string
	Email    string
	PwHash   string
}

type Message struct {
	Message_id uint `gorm:"primaryKey;autoIncrement"`
	Author_id  uint
	Text       string
	Pub_date   string
	Flagged    int
}

type Follower struct {
	WhoID  uint
	WhomID uint
}
