package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/providers"
	"github.com/varakornpz/utils"
)

var mqttClient mqtt.Client


func CommandHandler(c fiber.Ctx) error {
	var payload struct {
		Action string `json:"action"`
		CaneID	string `json:"cane_id"`
	}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "Invalid request" , "ok" : false})
	}
	userUUID , uuidErr := utils.GetUUIDByContext(c)
	if uuidErr != nil {
		log.Warn().Msgf("Cane get user uuid." , )
		return  c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"msg" : "Internal server error." ,
			"ok" : false ,
		})
	}
	isContain , containErr := utils.IsThisUserContainThisCane(c , payload.CaneID)
	if containErr != nil {
		log.Warn().Msgf("Cane contain check for uuid %s failed." , userUUID)
		return  c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"msg" : "Internal server error." ,
			"ok" : false ,
		})
	}

	if !isContain {
		log.Warn().Msgf("Cane contain check for uuid %s failed." , userUUID)
		return  c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg" : "Forbidden." ,
			"ok" : false ,
		})
	}

	topic := providers.AppConf.MQTTTopicPrefix +"/" + payload.CaneID + "/action"

	text := payload.Action

	token := mqttClient.Publish(topic , 0  ,false , text)
	token.Wait()

	return c.JSON(fiber.Map{
		"msg" : "Success." ,
		"action" : payload.Action ,
		"target" : topic ,
		"ok" : true ,
	})
}
