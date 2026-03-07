package myline

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/providers"
)


var myBot *messaging_api.MessagingApiAPI

func InitLine(){
	bot , err := messaging_api.NewMessagingApiAPI(providers.AppConf.LINECHANNELACC)
	if err != nil {
		log.Error().Msgf("Create line msg bot instance error.\n %s" , err)
	}
	myBot = bot
}


func BroadcastFallAlertWithoutLocation(caneID string){
if myBot == nil {
		log.Error().Msg("Cannot broadcast: myBot is nil")
		return
	}
	flexJSON := fmt.Sprintf(`{
		"type": "bubble",
		"size": "mega",
		"body": {
			"type": "box",
			"layout": "vertical",
			"contents": [
				{
					"type": "text",
					"text": "🚨 Alert!",
					"weight": "bold",
					"color": "#FF0000",
					"size": "xl"
				},
				{
					"type": "text",
					"text": "User may fall",
					"size": "md",
					"margin": "md",
					"wrap": true
				},
				{
					"type": "text",
					"text": "Cane ID : %s",
					"color": "#888888",
					"size": "sm",
					"margin": "sm"
				}
			]
		},
		"footer": {
					"type": "box",
					"layout": "vertical",
					"contents": [
						{
							"type": "text",
							"text": "⚠️ Latest location not found",
							"color": "#E53935",
							"weight": "bold",
							"size": "sm",
							"align": "center",
							"wrap": true
						}
					]
				}
}`, caneID) 


	container, err := messaging_api.UnmarshalFlexContainer([]byte(flexJSON))
	if err != nil {
		log.Error().Msgf("Parse Flex JSON error: %v", err)
		return
	}


	flexMsg := messaging_api.FlexMessage{
		AltText:  fmt.Sprintf("⚠️ Alert: fall detected (Cane ID : %s)", caneID),
		Contents: container,
	}

	_, err = myBot.Broadcast(
		&messaging_api.BroadcastRequest{
			Messages: []messaging_api.MessageInterface{flexMsg},
		},
		"",
	)

	if err != nil {
		log.Error().Msgf("Broadcast Flex error: %v", err)
	}
}

func BroadcastFallAlert(caneID string, lat float64, lng float64){
if myBot == nil {
		log.Error().Msg("Cannot broadcast: myBot is nil")
		return
	}

	flexJSON := fmt.Sprintf(`{
		"type": "bubble",
		"size": "mega",
		"body": {
			"type": "box",
			"layout": "vertical",
			"contents": [
				{
					"type": "text",
					"text": "🚨 Alert!",
					"weight": "bold",
					"color": "#FF0000",
					"size": "xl"
				},
				{
					"type": "text",
					"text": "User may fall",
					"size": "md",
					"margin": "md",
					"wrap": true
				},
				{
					"type": "text",
					"text": "Cane ID: %s",
					"color": "#888888",
					"size": "sm",
					"margin": "sm"
				}
			]
		},
		"footer": {
			"type": "box",
			"layout": "vertical",
			"contents": [
				{
					"type": "button",
					"action": {
						"type": "uri",
						"label": "📍 See location",
						"uri": "https://maps.google.com/?q=%f,%f"
					},
					"style": "primary",
					"color": "#E53935"
				}
			]
		}
}`, caneID, lat, lng)

	container, err := messaging_api.UnmarshalFlexContainer([]byte(flexJSON))
	if err != nil {
		log.Error().Msgf("Parse Flex JSON error: %v", err)
		return
	}

	flexMsg := messaging_api.FlexMessage{
		AltText:  fmt.Sprintf("⚠️ Alert: Fall detected (Cane ID %s)", caneID),
		Contents: container,
	}


	_, err = myBot.Broadcast(
		&messaging_api.BroadcastRequest{
			Messages: []messaging_api.MessageInterface{flexMsg},
		},
		"",
	)

	if err != nil {
		log.Error().Msgf("Broadcast Flex error: %v", err)
	}
}


func BroadCastToLine(msg string){
	_, err := myBot.Broadcast(
			&messaging_api.BroadcastRequest{
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: msg,
					},
				},
			},
			"",
		)

		if err != nil {
		log.Error().Msgf("Broadcast error: %v", err)
	} else {
		log.Info().Msg("Broadcast success!")
	}
}
