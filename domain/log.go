package domain

import "time"

type Log struct {
	UserId    string    `bson:"userid"`
	TimeStamp time.Time `bson:"timestamp"`
	Info      string    `bson:"info"`
	Function  string    `bson:"function"`
}
