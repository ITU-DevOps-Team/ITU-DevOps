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
	Pub_date   int
	Flagged    int
}

type Follower struct {
	WhoID  uint
	WhomID uint
}

//STRUCTS FOR VIEW
type ViewPost struct {
	Username   string
	Message_id uint
	Author_id  uint
	Text       string
	Pub_date   string
	Flagged    int
}

type ViewContent struct {
	SignedIn bool
	Posts []ViewPost
	User User
	Error bool
	Success bool
	ErrorMessage string
	SuccessMessage string
	SameUser bool //Personal timeline
	Username string
	AlreadyFollowing bool
}