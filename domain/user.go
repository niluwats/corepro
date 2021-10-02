package domain

type User struct {
	Id             string `bson:"_id,omitempty"`
	Name           string `bson:"name"`
	Email          string `bson:"email"`
	Password       string `bson:"password"`
	EmailVerified  bool   `bson:"email_verified"`
	MobileVerified bool   `bson:"mobile_verified"`
	Activated      bool   `bson:"activated"`
}

type FileUploadStatus struct {
	BirthRegImg      bool
	NICImg           bool
	PostalIdImg      bool
	DrivingLicensImg bool
	PassportImg      bool
}