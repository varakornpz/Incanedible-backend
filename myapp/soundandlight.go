package myapp

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/mqtt"
	"github.com/varakornpz/utils"
)

func GetSoundAndLight(c *websocket.Conn) {
    rawUUID := c.Locals("uuid")
    caneID := c.Query("cane_id")
    if caneID == "" {
        c.WriteMessage(websocket.TextMessage, []byte("Error: cane id is required"))
        c.Close()
        return
    }

    realUUID , uuidErr := utils.AnyToUUID(rawUUID)
    if uuidErr != nil {

        c.WriteJSON( fiber.Map{
            "msg" : "Internal convert uuid error" ,
            "ok" : false ,
        })
        c.Close()
        return
    }

    isContain , containErr := utils.IsThisUserUUIDContainThisCane(realUUID, caneID)
    if containErr != nil {
        c.WriteJSON( fiber.Map{
            "msg" : "Internal contain cane check error." ,
            "ok" : false ,
        })
        c.Close()
        return
    }

    if !isContain {
        c.WriteJSON( fiber.Map{
            "msg" : "Forbidden , this cane is not yours" ,
            "ok" : false ,
        })
        c.Close()
        return
    }


    topic := "b6810503897/" + caneID + "/soundandlight"


    msgChan := make(chan string, 50) 
    
    mqtt.Handler.AddClient(topic, msgChan)
    log.Info().Str("topic", topic).Msg("Start Subscribing...")

    defer func() {
        mqtt.Handler.RemoveClient(topic, msgChan)
        close(msgChan)
        c.Close()
        log.Info().Msg("Client Disconnected, Cleaned up resources")
    }()

    c.WriteMessage(websocket.TextMessage, []byte("Connected successfully to topic: "+topic))

    go func() {
        for {
            if _, _, err := c.ReadMessage(); err != nil {
                log.Debug().Msg("Client closed connection (ReadMessage Error)")
                break
            }
        }
    }()

    for {
        msgPayload, ok := <-msgChan
        
        if !ok {
            break
        }

        if err := c.WriteMessage(websocket.TextMessage, []byte(msgPayload)); err != nil {
            log.Error().Err(err).Msg("Write Error (Client likely disconnected)")
            break
        }
    }
}