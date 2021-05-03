package main

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getTimeFromTimestamp(timestamp int) string {
	timestampAsString := strconv.Itoa(timestamp)
	i, err := strconv.ParseInt(timestampAsString, 10, 64)
	checkErr(err)
	tm := time.Unix(i, 0)
	return tm.String()
}

func CheckUsernameExists(username string, db *gorm.DB) bool {
	user, err := GetUserByUsername(username, db)
	if err != nil {
		return false
	} else {
		return username == user.Username
	}
}
