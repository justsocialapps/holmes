package publisher

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/Shopify/sarama.v1"

	"github.com/justsocialapps/holmes/tracker"
	"github.com/justsocialapps/justlib"
)

// Publish creates a Kafka producer and provides it with every TrackingObject
// received via the given trackingChannel. If the producer cannot be started
// the program is quit.
func Publish(trackingChannel <-chan *tracker.TrackingObject, kafkaHost *string, kafkaTopic string) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Version = sarama.V0_10_0_0
	var producer sarama.AsyncProducer
	err := justlib.Try(0, 10*time.Second, func() error {
		var err error
		producer, err = sarama.NewAsyncProducer([]string{*kafkaHost}, kafkaConfig)
		if err != nil {
			log.Printf("Connection attempt to Kafka broker failed; trying again...\n")
		}
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
			log.Printf("successfully delivered msg to offset %d and partition %d\n", successMsg.Offset, successMsg.Partition)
		case errorMsg := <-producer.Errors():
			log.Printf("error delivering message: %v\n", *errorMsg)
		case object = <-trackingChannel:
			stringifiedObject, err := json.Marshal(object)
			if err == nil {
				log.Printf("publishing %s", stringifiedObject)
				producer.Input() <- &sarama.ProducerMessage{Topic: kafkaTopic, Key: nil, Value: sarama.ByteEncoder(stringifiedObject), Timestamp: time.Now()}
			}
		}
	}
}
