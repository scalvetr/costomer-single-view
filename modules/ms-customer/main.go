package main

import (
	"encoding/json"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mycujoo/go-kafka-avro"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	keySchemaFile := flag.String("key-schema-file", "customer-key.avro", "AVRO key schema file")
	valueSchemaFile := flag.String("value-schema-file", "customer-value.avro", "AVRO value schema file")
	dataFile := flag.String("data-file", "data.json", "JSON data file")
	bootstrapServers := getEnv("BOOTSTRAP_SERVERS", "localhost:9092")
	schemaRegistryUrl := getEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081")
	topicName := getEnv("TOPIC_NAME", "event.customer.entity")
	flag.Parse()
	log.Printf("keySchemaFile: %v\n", *keySchemaFile)
	log.Printf("valueSchemaFile: %v\n", *valueSchemaFile)
	log.Printf("dataFile: %v\n", *dataFile)
	log.Printf("bootstrapServers: %v\n", bootstrapServers)
	log.Printf("schemaRegistryUrl: %v\n", schemaRegistryUrl)
	log.Printf("topicName: %v\n", topicName)

	log.Printf("readSchema(#{*keySchemaFile})\n")
	keySchema := readSchema(*keySchemaFile)
	log.Printf("readSchema(#{*valueSchemaFile})\n")
	valueSchema := readSchema(*valueSchemaFile)

	log.Printf("kafka.NewProducer\n")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	log.Printf("kafkaavro.NewCachedSchemaRegistryClient\n")
	srClient, err := kafkaavro.NewCachedSchemaRegistryClient(schemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	log.Printf("kafkaavro.NewProducer\n")
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

	log.Printf("readData\n")
	data := readData(*dataFile)
	for _, item := range data {
		log.Printf("avroProducer.Produce\n")
		avroProducer.Produce(item.Id, item, nil)
		log.Printf("item sent: %v\n", item)
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func readSchema(schemaName string) string {
	avroSchemaBytes, err := ioutil.ReadFile(schemaName)
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string and print to screen
	avroSchema := string(avroSchemaBytes)
	//fmt.Println(avroSchema)
	return avroSchema
}
func readData(fileName string) []CustomerStruct {
	log.Printf("readData: ioutil.ReadFile\n")
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var customer []CustomerStruct

	log.Printf("readData: json.Unmarshal\n")
	err = json.Unmarshal(fileBytes, &customer)
	if err != nil {
		log.Fatal(err)
	}
	return customer
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
	Telephones []TelephoneStruct
	Addresses  []AddressStruct
}
