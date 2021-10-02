package dbhandler

import (
	"context"
	"core/domain"
	"core/errs"
	"core/logger"
	"core/utils"
	"log"

	"time"

	"gopkg.in/mgo.v2/bson"
)

//	CreateUser gets user struct and query the given user is already exists in DB. If exists return an error. If not exists create a new user in DB. returns a user struct and and error
func CreateUser(user domain.User) (*domain.User, *errs.AppError) {
	col := Db.Collection("support")
	u := &domain.User{}

	ifUser := FindIfEmailExists(user.Email)
	if ifUser == false {
		id := bson.NewObjectId()
		user.Id = id.Hex()
		_, err := col.InsertOne(context.TODO(), &user)
		if err != nil {
			SaveUserLog(user.Id, logger.Error("failed to insert user into DB"))
			return nil, errs.NewUnexpectedError(err.Error())
		}

		code := utils.GenerateCode("verify_email")
		now := time.Now()
		timeOut := now.Add(time.Minute * 45)
		err1 := SaveEmailVerificationCode(code, user.Email, timeOut)
		if err1 != nil {
			SaveUserLog(user.Id, logger.Error(err1.Message))
			return nil, errs.NewUnexpectedError(err1.Message)
		}

		err2 := utils.SendEmail(code, user.Email, "userreg")
		if err2 != nil {
			SaveUserLog(user.Id, logger.Error(err2.Message))
			return nil, errs.NewUnexpectedError(err2.Message)
		}
	} else {
		SaveUserLog(user.Id, logger.Error(errs.Exist_user))
		return nil, errs.NewUnexpectedError(errs.Exist_user)
	}
	SaveUserLog(user.Id, logger.Info("created new user"))
	return u, nil
}

//	SaveUserLog gets user ID and a struct type of *Log. Insert log into DB.
func SaveUserLog(userid string, lg *domain.Log) {
	lg.UserId = userid
	col := Db.Collection("logger")
	_, err := col.InsertOne(context.TODO(), &lg)
	if err != nil {
		log.Fatal(err)
	}
}
