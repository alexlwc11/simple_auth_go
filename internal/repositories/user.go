package repositories

import (
	"github.com/alexlwc11/simple_auth_go/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateWithDeviceUUID(deviceUUID string) (*models.User, error)
	FindByDeviceUUID(deviceUUID string) (*models.User, error)
	// CreateWithEmailPassword(email string, password string) (*models.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (ur *UserRepositoryImpl) CreateWithDeviceUUID(deviceUUID string) (*models.User, error) {
	user := models.User{
		DeviceUUID: deviceUUID,
	}
	error := ur.DB.Create(&user).Error
	if error != nil {
		return &models.User{}, error
	}

	return &user, nil
}

func (ur *UserRepositoryImpl) FindByDeviceUUID(deviceUUID string) (*models.User, error) {
	var user models.User
	err := ur.DB.Where("device_uuid = ?", deviceUUID).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}
