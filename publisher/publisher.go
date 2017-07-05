package publisher

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/Shopify/sarama.v1"

	"github.com/justsocialapps/holmes/tracker"
	"github.com/justsocialapps/justlib"
)

func Publish(trackingChannel <-chan *tracker.TrackingObject, kafkaHost *string, kafkaTopic string) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Version = sarama.V0_10_0_0
	var producer sarama.AsyncProducer
	err := justlib.Try(7, 10*time.Second, func() error {
		var err error
		producer, err = sarama.NewAsyncProducer([]string{*kafkaHost}, kafkaConfig)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Kafka producer up and running with broker %s", *kafkaHost)

	var object *tracker.TrackingObject
	for {
		select {
		case successMsg := <-producer.Successes():
			log.Printf("successfully delivered msg: %v\n", successMsg.Offset)
		case errorMsg := <-producer.Errors():
			log.Printf("error delivering message: %v\n", *errorMsg)
		case object = <-trackingChannel:
			stringifiedObject, err := json.Marshal(object)
			if err == nil {
				log.Printf("publishing %s", stringifiedObject)
				producer.Input() <- &sarama.ProducerMessage{Topic: kafkaTopic, Key: sarama.StringEncoder("key"), Value: sarama.ByteEncoder(stringifiedObject), Timestamp: time.Now()}
			}
		}
	}
}
