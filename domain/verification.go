package domain

import (
	"time"
)

type EmailVerification struct {
	Email    string    `bson:"email"`
	Code     string    `bson:"code"`
	TimeOut  time.Time `bson:"timeout"`
	OldEmail string    `bson:"oldemail"`
}

type MobileVerification struct {
	UserId       string    `bson:"user_id"`
	MobileNumber string    `bson:"mobile_number"`
	Code         string    `bson:"code"`
	TimeOut      time.Time `bson:"timeout"`
}
