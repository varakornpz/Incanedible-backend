package myapp

import (
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/mygorm"
	"github.com/varakornpz/utils"

	"github.com/gofiber/fiber/v3"
)


func GetUserData(c fiber.Ctx) error{
	userUUID , uuidErr := utils.GetUUIDByContext(c)
	if uuidErr != nil {
		log.Fatal().Msg(uuidErr.Error())
	}

	user , getErr := mygorm.GetUserByUUID(userUUID)
	if getErr != nil {
		log.Error().Msg("get user error in myapp.go")
		return  c.SendStatus(fiber.ErrInternalServerError.Code)
	}


	return c.JSON(fiber.Map{
		"email" : user.Email ,
		"name" : user.Name ,
		"profile_pic" : user.ProfilePic ,
		"canes" : user.RegisteredCanes ,
	})
}