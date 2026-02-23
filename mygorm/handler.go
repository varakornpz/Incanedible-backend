package mygorm

import (
	"github.com/google/uuid"
	"github.com/varakornpz/models"
)


func GetUserByEmail(email string) (models.User , error){
	var user models.User
	result := DB.First(&user ,"email = ?" , email)

	return user , result.Error
}

func GetUserByUUID(uuid uuid.UUID) (models.User , error){
	var user models.User
	result := DB.First(&user , "uuid = ?" , uuid)

	return user , result.Error
}

func GetCanesByUUID(uuid uuid.UUID)(models.RegisteredCanes , error){
	var result struct {
        Canes models.RegisteredCanes `gorm:"column:registered_canes"` 
    }

	dbResult := DB.Model(&models.User{}).
        Select("registered_canes").
        Where("uuid = ?", uuid).
        First(&result)
	return  result.Canes , dbResult.Error
}


func PutNewUser(user *models.User) error {
	result := DB.Create(user)
	return  result.Error
}