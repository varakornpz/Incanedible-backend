//Gemini 90%

package mqtt

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/providers"
	myredis "github.com/varakornpz/redis"
)

var (
    Client  mqtt.Client
    Handler *MQTTSubscibeHandler
	SoundLightChan = make(chan string, 100)
)

func InitMQTT() {
	opts := mqtt.NewClientOptions()

	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(15 * time.Second)

	opts.AddBroker(providers.AppConf.MQTTBroker + ":1883")
	opts.SetUsername(providers.AppConf.MQTTUsername)
	opts.SetPassword(providers.AppConf.MQTTPassword)
	opts.SetClientID(providers.AppConf.MQTTClientID)

	Handler = &MQTTSubscibeHandler{
		clients: make(map[string]map[chan string]bool),
	}

	opts.OnConnect = func(c mqtt.Client) {
		log.Info().Msg("Connected to MQTT Broker. Subscribing to topics...")

		locationTopic := providers.AppConf.MQTTTopicPrefix + "/+/location"
		if token := c.Subscribe(locationTopic, providers.AppConf.MQTTQos, Handler.handleMQTTMessage); token.Wait() && token.Error() != nil {
			log.Error().Msgf("ไม่สามารถ Subscribe Location ได้: %v", token.Error())
		} else {
			log.Info().Msgf("เริ่มดักฟัง Location ตลอดเวลาที่ Topic: %s", locationTopic)
		}

		soundLightTopic := providers.AppConf.MQTTTopicPrefix + "/+/soundandlight"
		if token := c.Subscribe(soundLightTopic, providers.AppConf.MQTTQos, Handler.handleMQTTMessage); token.Wait() && token.Error() != nil {
			log.Error().Msgf("ไม่สามารถ Subscribe Sound & Light ได้: %v", token.Error())
		} else {
			log.Info().Msgf("เริ่มดักฟัง Sound & Light ตลอดเวลาที่ Topic: %s", soundLightTopic)
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Warn().Msgf("Connect lost: %v", err)
	}

	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("Failed to connect MQTT: %s", token.Error())
	}

	Handler.mqttCli = Client

	go SubscribeFall()

	log.Info().Msg("MQTT Subscription Handler Initialized Successfully")
}

type MQTTSubscibeHandler struct {
	mu      sync.RWMutex
	clients map[string]map[chan string]bool
	mqttCli mqtt.Client
}


func (sm *MQTTSubscibeHandler ) AddClient(topic string, ch chan string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.clients[topic] == nil {
		sm.clients[topic] = make(map[chan string]bool)
		log.Printf("Subscribing to new topic on broker: %s\n", topic)
		if token := sm.mqttCli.Subscribe(topic, providers.AppConf.MQTTQos, sm.handleMQTTMessage); token.Wait() && token.Error() != nil {
             log.Error().Msgf("Error subscribing: %v", token.Error())
        }
	}
	sm.clients[topic][ch] = true
}

func (sm *MQTTSubscibeHandler) RemoveClient(topic string, ch chan string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.clients[topic] != nil {
		delete(sm.clients[topic], ch)

		if len(sm.clients[topic]) == 0 {
			log.Printf("Unsubscribing from topic on broker: %s\n", topic)
			delete(sm.clients, topic)
			sm.mqttCli.Unsubscribe(topic)
		}
	}
}



func (sm *MQTTSubscibeHandler) handleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())



	if strings.HasSuffix(topic , "/location"){
		parts := strings.Split(topic, "/")


		if len(parts) < 3 {
			log.Warn().Msgf("Topic error recieve incorect format : %s" , topic)
			return
		}else{

			caneID := parts[1]

			var jsonPayload myredis.LatestCaneLocation

			jsonErr := json.Unmarshal(msg.Payload() , &jsonPayload)

			if jsonErr != nil {
				log.Error().Msgf("Cant Unmarshal payload,format may not correct : %v" , payload)
			}else{
				go func(){
					err := myredis.PutCaneAddress(caneID , jsonPayload)

					if err != nil {
						log.Error().Msgf("Redis background error : %v" , err)
					}
				}()
			}
		}
	}

	

	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for ch := range sm.clients[topic] {
		select {
		case ch <- payload:
		default:

		}
	}
}
