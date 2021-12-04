package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mycujoo/go-kafka-avro"
	"io/ioutil"
	"log"
	"os"
)

func readSchema(schemaName string) string {
	avroSchemaBytes, err := ioutil.ReadFile(schemaName)
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string and print to screen
	avroSchema := string(avroSchemaBytes)
	fmt.Println(avroSchema)
	return avroSchema
}
func main() {
	keySchema := readSchema("customer-key.avro")
	valueSchema := readSchema("customer-value.avro")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS")})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	srClient, err := kafkaavro.NewCachedSchemaRegistryClient(os.Getenv("SCHEMA_REGISTRY_URL"))
	if err != nil {
		panic(err)
	}

	topicName := os.Getenv("TOPIC_NAME")

	avroProducer, err := kafkaavro.NewProducer(kafkaavro.ProducerConfig{
		TopicName:            topicName,
		KeySchema:            keySchema,
		ValueSchema:          valueSchema,
		Producer:             p,
		SchemaRegistryClient: srClient,
	})
	if err != nil {
		panic(err)
	}
	defer avroProducer.Close()

	avroProducer.Produce("key", "value", nil)

	data := readData("data.json")
	for _, item := range data {
		avroProducer.Produce(item.Id, item, nil)
	}

}

func readData(location string) []CustomerStruct {

}

type AddressStruct struct {
	Street  string
	Number  string
	City    string
	Country string
	ZipCode string
	Default bool
}
type TelephoneStruct struct {
	Number  string
	Primary bool
}
type CustomerStruct struct {
	Id         string
	Name       string
	Surname    string
	Email      string
	Telephones TelephoneStruct
	Addresses  AddressStruct
}
