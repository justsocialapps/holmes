package publisher

import (
	"encoding/json"
	"log"

	"github.com/shopify/sarama"

	"github.com/justsocialapps/holmes/models"
)

func Publish(trackingChannel <-chan *models.TrackingObject, kafkaHost *string, kafkaTopic string) error {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Errors = false
	producer, err := sarama.NewAsyncProducer([]string{*kafkaHost}, kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}

	var object *models.TrackingObject
	for {
		object = <-trackingChannel
		stringifiedObject, err := json.Marshal(object)
		if err == nil {
			log.Printf("publishing %s", stringifiedObject)
			producer.Input() <- &sarama.ProducerMessage{Topic: kafkaTopic, Key: sarama.StringEncoder("key"), Value: sarama.ByteEncoder(stringifiedObject)}
		}
	}
}
