package main

import (
	"os"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/auth"
	myline "github.com/varakornpz/line"
	"github.com/varakornpz/mqtt"
	"github.com/varakornpz/myapp"
	"github.com/varakornpz/mygorm"
	"github.com/varakornpz/providers"
	myredis "github.com/varakornpz/redis"
	"github.com/varakornpz/utils"
)



func main(){
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	//Provider must be on the top
	providers.InitAppConf()
	myredis.InitRedis()
	myline.InitLine()
	auth.InitGoogleAuthConf()

	app := fiber.New()
	InitCORSConf(app)
	mqtt.InitMQTT()
	app.Get("/" , func (c fiber.Ctx) error  {
		return c.SendString("huh , wtf is this place.")
	})

	authRoute := app.Group("/auth")

	authRoute.Get("/google/callback" , auth.GoogleAuthCallBack)
	authRoute.Get("/google/signin" , auth.GoogleAuthSignin)

	mainAppRoute := app.Group("/app")
	mainAppRoute.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(providers.AppConf.JWTSecret)},
		Extractor : extractors.FromCookie("access_token") ,
		ErrorHandler: func (c fiber.Ctx , err error) error {
			if err.Error() == "value not found" {
				return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
					"msg" : "Token not found" ,
					"ok" : false ,
				})
			}
			log.Error().Msgf("JWT Authentication Failed: %v", err)
			return c.SendStatus(fiber.ErrUnauthorized.Code)
		},
	}))
	mainAppRoute.Get("/hi" , func (c fiber.Ctx) error  {
		return c.JSON(fiber.Map{
			"msg" : "hi" ,
		})
	})

	mainAppRoute.Post("/addcane" , myapp.AddCane)

	mainAppRoute.Post("/action" , mqtt.CommandHandler)

	mainAppSocketRoute := mainAppRoute.Group("/socket")

	mainAppSocketRoute.Use("/" , func(c fiber.Ctx) error {
		uuid , uuidErr := utils.GetUUIDByContext(c)
		if uuidErr != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"msg" : "Cant get uuid by user context" ,
				"ok" : false ,
			})
		}
		c.Locals("uuid" , uuid)

		return c.Next()
	})


	mainAppSocketRoute.Get("/getlocation" , websocket.New(myapp.GetLocation))

	mainAppRoute.Get("/me" , myapp.GetUserData)


	mygorm.InitDB()

	
	myline.BroadCastToLine("Hello init")

	app.Listen(":3334")
}