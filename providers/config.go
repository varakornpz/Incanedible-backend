package providers

import (
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)


type Config struct {
	JWTSecret		string	`validate:"required"`
	DBDsn			string	`validate:"required"`
	GGClientSecret	string	`validate:"required"`
	GGClientID		string	`validate:"required"`
	GGRedirectUrl	string	`validate:"required"`
	GGAfterSigninRedirect	string	`validate:"required"`
	FrontendErrPage		string	`validate:"required"`
	COOKIEDomain	string	`validate:"required"`
	MQTTQos		byte		`validate:"required"`
	MQTTBroker	string	`validate:"required"`
	MQTTTopicPrefix		string		`validate:"required"`
	MQTTUsername	string		`validate:"required"`
	MQTTPassword	string		`validate:"required"`
	MQTTClientID	string		`validate:"required"`
}

var AppConf	*Config

func InitAppConf(){
	dotEnvErr := godotenv.Load()
    if dotEnvErr != nil {
		log.Info().Msg("No .env file found, using system environment variables")
    }

	mqttQosStr := os.Getenv("MQTT_QOS")
	var mqttQos byte
	mqttQosInted, qosErr := strconv.Atoi(mqttQosStr)
	if qosErr != nil {
		log.Fatal().Err(qosErr).Msg("MQTT_QOS must be a valid integer")
	}
	mqttQos = byte(mqttQosInted)

	AppConf = &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		DBDsn: os.Getenv("DB_DSN"),
		GGClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GGClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		GGRedirectUrl: os.Getenv("GOOGLE_REDIRECT_URL"),
		GGAfterSigninRedirect: os.Getenv("GOOGLE_AFTER_SIGNIN_REDIRECT"),
		FrontendErrPage: os.Getenv("FE_ERROR_PAGE") ,
		COOKIEDomain: os.Getenv("COOKIE_DOMAIN"),
		MQTTQos: mqttQos ,
		MQTTBroker: os.Getenv("MQTT_BROKER"),
		MQTTTopicPrefix: os.Getenv("MQTT_TOPIC_PREFIX"),
		MQTTUsername: os.Getenv("MQTT_USERNAME"),
		MQTTPassword: os.Getenv("MQTT_PASSWORD"),
		MQTTClientID: os.Getenv("MQTT_CLIENT_ID"),
		
	}

	v := reflect.ValueOf(*AppConf)
	t := v.Type()
	for i := 0 ; i < v.NumField() ; i++ {
		field := t.Field(i)
		fieldVal	:= v.Field(i)

		validateTag	:= field.Tag.Get("validate")
		
		if validateTag == "required" {
			if fieldVal.IsZero() {
				log.Fatal().Msgf("Missing required ENV: %s called by providers/config.go", field.Name)
			}
		}
	}

	log.Info().Msg("✅ Config loaded.")
}