package main

type User struct {
	UserID   uint
	Username string
	Email    string
	PwHash   string
}

type Message struct {
	Message_id uint
	Author_id  uint
	Text       string
	Pub_date   string
	Flagged    int
}
