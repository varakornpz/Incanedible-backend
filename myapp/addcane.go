package myapp

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rs/zerolog/log"

	// "github.com/varakornpz/models"
	"github.com/asaskevich/govalidator"
	"github.com/varakornpz/models"
	"github.com/varakornpz/mygorm"
	"github.com/varakornpz/utils"
)

type AddCaneReqBody struct {
	Name string `json:"name"`
	CaneID string `json:"cane_id"`
}

func AddCane(c fiber.Ctx) error {
	userUUID , uuidErr := utils.GetUUIDByContext(c)
	if uuidErr != nil {
		log.Fatal().Msg(uuidErr.Error())
	}

	user , getErr := mygorm.GetUserByUUID(userUUID)
	if getErr != nil {
		log.Error().Msg("get user error in myapp.go")
		return  c.SendStatus(fiber.ErrInternalServerError.Code)
	}

	var req AddCaneReqBody

	if err := c.Bind().JSON(&req) ; err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"msg" : "Body should be JSON" ,
			"ok" : false ,
		})
	}

	if(!govalidator.IsAlphanumeric(req.CaneID)){
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"msg" : "Cane id Can be mixed of number and eng alphabet." ,
			"ok" : false ,
		})
	}

	if len(req.CaneID) > 15 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"msg" : "Cane id cant longer than 15 digit." ,
			"ok" : false ,
		})
	}
	if len(req.CaneID) < 3 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"msg" : "Cane id length should more than 3 digit." ,
			"ok" : false ,
		})
	}




	newCane := models.RegisteredCane{
		Name : "" ,
		CaneID: req.CaneID,
	}
	if len(user.RegisteredCanes) >= 5{
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"msg" : "You hit a cane limit,And you cant delete it by yourself :)" ,
			"ok" : false ,
		})
	}

	

	user.RegisteredCanes = append(user.RegisteredCanes, newCane)

	mygorm.DB.Save(user)

	return c.JSON(fiber.Map{
        "msg": "Cane added successfully",
        "user": user,
		"ok" : true ,
    })
	

}
