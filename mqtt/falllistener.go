package mqtt

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	myline "github.com/varakornpz/line"
	"github.com/varakornpz/providers"
	myredis "github.com/varakornpz/redis"
)

type FallPayload struct {
	DeviceID	string	`json:"deviceID"`
}


func SubscribeFall(){
	log.Debug().Msgf("Start listening to %s/+/fall" , providers.AppConf.MQTTTopicPrefix) ;


	topicString := providers.AppConf.MQTTTopicPrefix + "/+/fall"
	Client.Subscribe(topicString , providers.AppConf.MQTTQos , func(client mqtt.Client, msg mqtt.Message){
		topic := msg.Topic()

		if strings.HasSuffix(topic , "/fall"){
		parts := strings.Split(topic, "/")


		if len(parts) < 3 {
			log.Warn().Msgf("Topic error recieve incorect format : %s" , topic)
			return
		}else{

			caneID := parts[1]

			var jsonPayload FallPayload


			// var jsonLocation myredis.LatestCaneLocation

			jsonErr := json.Unmarshal(msg.Payload() , &jsonPayload)

			if jsonErr != nil {
				log.Error().Msgf("Cant Unmarshal payload,format may not correct : %v" , msg.Payload())
			}else{
				go func(cID string){
					location , err := myredis.GetLatestLocation(cID)
					if err != nil {
						log.Error().Msgf("Redis background error : %v" , err)
						myline.BroadcastFallAlertWithoutLocation(cID)
					}else{
						myline.BroadcastFallAlert(cID , location.Lat  , location.Lng)
					}

				}(caneID)
			}
		}
	}
		log.Debug().Msgf("Recieve fall from %s" , msg.Topic())
	})

}