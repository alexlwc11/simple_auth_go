package models

// TODO support other sign up methods e.g. email & password
type User struct {
	BaseModel
	DeviceUUID string `gorm:"uniqueIndex;size:36;not null"`
	// Credentials
}

type Credentials struct {
	Email    string `gorm:"uniqueIndex"`
	Password string
}
