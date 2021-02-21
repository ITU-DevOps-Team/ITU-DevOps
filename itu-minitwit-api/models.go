package main

type User struct {
	UserID   uint `gorm:"primaryKey;autoIncrement"`
	Username string
	Email    string
	PwHash   string
}
type Follower struct {
	WhoID  uint
	WhomID uint
}

type Message struct {
	Message_id uint `gorm:"primaryKey;autoIncrement"`
	Author_id  uint
	Text       string
	Pub_date   int
	Flagged    int
}

type User_ struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}

type Message_ struct {
	Content string `json:"content"`
}

type Response struct {
	Status    int    `json:"status"`
	Error_msg string `json:"error_msg"`
}

type Latest struct {
	Latest int `json:"latest"`
}
