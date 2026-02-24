//Gemini 90%

package mqtt

import (

	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/providers"
)

var (
    Client  mqtt.Client
    Handler *MQTTSubscibeHandler
)

func InitMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(providers.AppConf.MQTTBroker + ":1883")
	opts.SetUsername(providers.AppConf.MQTTUsername)
	opts.SetPassword(providers.AppConf.MQTTPassword)
	opts.SetClientID(providers.AppConf.MQTTClientID)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Info().Msg("Connected to MQTT Broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Warn().Msgf("Connect lost: %v", err)
	}

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("%s" , token.Error())
	}

	Handler = &MQTTSubscibeHandler{
        clients: make(map[string]map[chan string]bool),
        mqttCli: mqttClient,
    }

	log.Info().Msg("MQTT Subscription Handler Initialized")
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

	sm.mu.RLock()
	defer sm.mu.RUnlock()




	for ch := range sm.clients[topic] {
		select {
		case ch <- payload:
		default:

		}
	}
}
