package dbhandler

import (
	"context"
	"core/domain"
	"core/errs"
	"core/logger"
	"core/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//	Login gets email and password received by service and queries the email and password seperately. If credentials are correct returns the queried user as a struct and an error.
func Login(email, password string) (*domain.User, *errs.AppError) {
	user := Db.Collection("support")
	cm := domain.User{}
	err := user.FindOne(context.TODO(), bson.M{"email": email}).Decode(&cm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			SaveUserLog("unknown", logger.Error(errs.Wrong_em))
			return nil, errs.NewNotFoundError(errs.Wrong_em)
		}
	}
	hashedPassword := cm.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		SaveUserLog(cm.Id, logger.Error(errs.Wrong_pw))
		return nil, errs.NewAuthenticationError(errs.Wrong_pw)
	}
	if cm.EmailVerified == false {
		SaveUserLog(cm.Id, logger.Error(errs.Unverified_em))
		return nil, errs.NewUnexpectedError(errs.Unverified_em)
	}
	if cm.Activated == false {
		SaveUserLog(cm.Id, logger.Error(errs.Deactivated_acc))
		return nil, errs.NewUnexpectedError(errs.Deactivated_acc)
	}
	SaveUserLog(cm.Id, logger.Info("logged in "))
	return &cm, nil
}

//	FindIfEmailExists query the given email and returns if it is available in DB or not
func FindIfEmailExists(email string) bool {
	var user domain.User
	col := Db.Collection("support")
	err := col.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return false
	} else {
		return true
	}
}

func SaveEmailVerificationCode(code, em string, timeOut time.Time) *errs.AppError {
	var ev domain.EmailVerification
	col := Db.Collection("email_verification")
	email_ver := domain.EmailVerification{
		Code:     code,
		TimeOut:  timeOut,
		Email:    em,
		OldEmail: "",
	}
	err := col.FindOne(context.TODO(), bson.M{"email": em}).Decode(&ev)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := col.InsertOne(context.TODO(), &email_ver)
			if err != nil {
				SaveUserLog("unknown", logger.Error("error while inserting email verification code into DB"))
				return errs.NewUnexpectedError(err.Error())
			}
		} else {
			SaveUserLog("unknown", logger.Error("error while querying email verification code into DB"))
			return errs.NewUnexpectedError(err.Error())
		}
	} else {
		_, err = col.UpdateOne(context.Background(), bson.M{"email": em}, bson.M{"$set": bson.M{"code": code, "timeout": timeOut}})
		if err != nil {
			SaveUserLog("unknown", logger.Error("error while updating email verification code into DB"))
			return errs.NewUnexpectedError(err.Error())
		}
	}
	SaveUserLog("unknown", logger.Info("saved email verification code"))
	return nil
}

//	VerifyEmail gets struct type of EmailVerification. Queries the user by email. Checks if given email in struct is matched with the given url hash. Updates user activation status in user collection and emailtoupdate in email verification collection. Returns an error.
func VerifyEmail(verify domain.EmailVerification) *errs.AppError {
	col := Db.Collection("users")
	cm := domain.User{}
	col.FindOne(context.TODO(), bson.M{"email": verify.Email}).Decode(&cm)
	ifSame, err := checkIfEqualEmailCode(verify.Code, verify.Email)
	if err != nil {
		return err
	}

	if ifSame == true {
		_, err := col.UpdateOne(context.Background(), bson.M{"email": verify.Email}, bson.M{"$set": bson.M{"activated": true, "email_verified": true}})
		if err != nil {
			SaveUserLog("unknown", logger.Error("error while updating email activation status"))
			return errs.NewUnexpectedError(err.Error())
		}
	} else {
		SaveUserLog("unknown", logger.Error("error verifying email"))
		return errs.NewUnexpectedError("error verifying email")
	}
	return nil
}

//	checkIfEqualEmailCode gets email verification hash and email. Queries them from email verification collection to check if url hash is same with the email as in collection. Returns bool and an error.
func checkIfEqualEmailCode(code, email string) (bool, *errs.AppError) {
	col := Db.Collection("email_verification")
	cm := domain.EmailVerification{}
	err := col.FindOne(context.TODO(), bson.M{"email": email}).Decode(&cm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			SaveUserLog("unknown", logger.Error("error while querying email"))
			logger.Error("error while " + err.Error())
		}
	}
	storedCode := cm.Code
	if storedCode != code {
		SaveUserLog("unknown", logger.Error("error while querying email"))
		return false, errs.NewUnexpectedError("verification code is not matching")
	}
	if time.Now().After(cm.TimeOut) {
		SaveUserLog("unknown", logger.Error("email verification timeout has expired"))
		return false, errs.NewUnexpectedError("timeout has expired")
	}
	SaveUserLog("unknown", logger.Info("email verification hash comparing succesfull"))
	return true, nil
}

func ResendActivationEmail(email string) *errs.AppError {
	user := Db.Collection("users")
	cm := domain.User{}
	user.FindOne(context.TODO(), bson.M{"email": email}).Decode(&cm)

	code := utils.GenerateCode("verify_email")
	now := time.Now()
	timeOut := now.Add(time.Minute * 45)
	err := SaveEmailVerificationCode(code, email, timeOut)
	if err != nil {
		SaveUserLog(cm.Id, logger.Error(err.Message))
		return err
	}
	err = utils.SendEmail(code, email, "userreg")
	if err != nil {
		SaveUserLog(cm.Id, logger.Error(err.Message))
		return errs.NewUnexpectedError(err.Message)
	}
	SaveUserLog(cm.Id, logger.Info("Resend activation email"))
	return nil
}

func RecoverEmail(email string) *errs.AppError {
	user := Db.Collection("users")
	cm := domain.User{}
	user.FindOne(context.TODO(), bson.M{"email": email}).Decode(&cm)
	now := time.Now()
	timeout := now.Add(time.Minute * 45)
	ifEx := FindIfEmailExists(email)

	if ifEx == true {
		code := utils.GenerateCode("recover_email")

		col := Db.Collection("email_verification")
		_, err := col.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"timeout": timeout, "code": code}})
		if err != nil {
			SaveUserLog(cm.Id, logger.Error("failed to recover email"))
			return errs.NewUnexpectedError(err.Error())
		}

		err2 := utils.SendEmail(code, email, "")
		if err2 != nil {
			SaveUserLog(cm.Id, logger.Error("error while sending email to "+email))
			return errs.NewUnexpectedError(err2.Message)
		}
	} else {
		SaveUserLog(cm.Id, logger.Error("failed to send recover email "))
		return errs.NewUnexpectedError("failed to send recover email")
	}
	return nil
}

func ResetPassword(evpw, url_email, email, newpw, confirmpw string) *errs.AppError {
	user := Db.Collection("users")
	cm := domain.User{}
	user.FindOne(context.TODO(), bson.M{"email": email}).Decode(&cm)

	if newpw != confirmpw {
		SaveUserLog(cm.Id, logger.Error("miss matched passwords "))
		return errs.NewUnexpectedError("passwords doesn't match")
	}
	if email != url_email {
		SaveUserLog(cm.Id, logger.Error("entered invalid email"))
		return errs.NewUnexpectedError("invalid email")
	}

	var emailverif domain.EmailVerification
	col := Db.Collection("email_verification")
	err := col.FindOne(context.TODO(), bson.M{"email": email, "code": evpw}).Decode(&emailverif)
	if err == mongo.ErrNoDocuments {
		SaveUserLog("unknown", logger.Error("incorrect hash in url"))
		return errs.NewUnexpectedError("incorrect hash in url")
	} else {
		if time.Now().After(emailverif.TimeOut) {
			SaveUserLog(cm.Id, logger.Error("reset password timeout"))
			return errs.NewUnexpectedError("timeout has expired")
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(confirmpw), 10)
			if err != nil {
				SaveUserLog(cm.Id, logger.Error("failed to reset password"))
				return errs.NewUnexpectedError(err.Error())
			}

			col = Db.Collection("users")
			_, err1 := col.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"password": hash}})
			if err1 != nil {
				SaveUserLog(cm.Id, logger.Error("failed to reset password"))
				return errs.NewUnexpectedError(err1.Error())
			}
		}
	}
	SaveUserLog(cm.Id, logger.Info("reset password"))
	return nil
}

func VerifyMobileNo(verify domain.MobileVerification) *errs.AppError {
	verify.Code = utils.Encode(verify.Code)
	ifSame, err := checkIfEqualSmsCode(verify.Code, verify.MobileNumber)
	if err != nil {
		SaveUserLog(verify.UserId, logger.Error(err.Message))
		return err
	}
	if ifSame == true {
		col := Db.Collection("users")
		_, err := col.UpdateOne(context.Background(), bson.M{"_id": verify.UserId}, bson.M{"$set": bson.M{"mobile_verified": true, "profile.contactdetails.mobile_number": verify.MobileNumber}})
		if err != nil {
			SaveUserLog(verify.UserId, logger.Error("error while updating mobile verification status -"))
			return errs.NewUnexpectedError(err.Error())
		} else {
			SaveUserLog(verify.UserId, logger.Info("mobile number verified"))
			return nil
		}
	} else {
		SaveUserLog(verify.UserId, logger.Error("error verifying mobile number"))
		return errs.NewUnexpectedError("error verifying mobile number")
	}
}

func SendMobileVerificationCode(mobile string, id string) *errs.AppError {
	if findIfMobileExists(mobile) == true {
		return errs.NewUnexpectedError("This mobile number is verified and taken already")
	}
	code := utils.GenerateCode("verify_mobile")
	now := time.Now()
	timeOut := now.Add(time.Minute * 45)
	err2 := utils.SendSms(mobile, code)
	if err2 != nil {
		SaveUserLog(id, logger.Error(err2.Message))
		return errs.NewUnexpectedError(err2.Message)
	}
	code_hashStr := utils.Encode(code)
	err1 := saveMobileVerificationCode(code_hashStr, mobile, id, timeOut)
	if err1 != nil {
		SaveUserLog(id, logger.Error(err1.Message))
		return errs.NewUnexpectedError("unexpected DB error")
	}
	SaveUserLog(id, logger.Info("mobile verification code sent"))
	return nil
}

func checkIfEqualSmsCode(code string, mobile string) (bool, *errs.AppError) {
	col := Db.Collection("mobile_verification")
	cm := domain.MobileVerification{}
	err := col.FindOne(context.TODO(), bson.M{"mobile_number": mobile}).Decode(&cm)
	if err != nil {
		SaveUserLog("unknown", logger.Error("error while querying mobile number"))
		return false, errs.NewUnexpectedError("phone number not found")
	}
	storedCode := cm.Code
	if storedCode != code {
		SaveUserLog("unknown", logger.Error("incorrect verification code"))
		return false, errs.NewUnexpectedError("verification code is not matching")
	}
	if time.Now().After(cm.TimeOut) {
		SaveUserLog("unknown", logger.Error("mobile verification timeout has expired"))
		return false, errs.NewUnexpectedError("timeout has expired")
	}
	SaveUserLog("unknown", logger.Info("mobile verification hash comparing succesfull"))
	return true, nil
}

func findIfMobileExists(mobile string) bool {
	var mb string
	col := Db.Collection("user")
	err := col.FindOne(context.TODO(), bson.M{"profile.contactdetails.mobile_number": mobile}).Decode(&mb)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func saveMobileVerificationCode(code, mobile, id string, timeOut time.Time) *errs.AppError {
	col := Db.Collection("mobile_verification")
	cm := domain.MobileVerification{
		Code:         code,
		MobileNumber: mobile,
		UserId:       id,
		TimeOut:      timeOut,
	}

	mobilev := domain.MobileVerification{}
	err := col.FindOne(context.TODO(), bson.M{"user_id": id}).Decode(&mobilev)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err1 := col.InsertOne(context.TODO(), &cm)
			if err1 != nil {
				SaveUserLog(id, logger.Error("error while inserting mobile verification code into DB "))
				return errs.NewUnexpectedError(err1.Error())
			}
		}
	} else {
		_, err := col.UpdateOne(context.Background(), bson.M{"user_id": id}, bson.M{"$set": bson.M{"code": code, "mobile_number": mobile, "timeout": timeOut}})
		if err != nil {
			SaveUserLog(id, logger.Error("error while updating mobile verification code "))
			return errs.NewUnexpectedError("error while updating mobile verification code")
		}
	}
	SaveUserLog(id, logger.Info("mobile verification code saved to DB"))
	return nil
}
