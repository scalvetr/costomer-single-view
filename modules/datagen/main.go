package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	keySchemaFile := flag.String("key-schema-file", "customer-key.avsc", "AVRO key schema file")
	valueSchemaFile := flag.String("value-schema-file", "customer-value.avsc", "AVRO value schema file")
	bootstrapServers := GetEnv("BOOTSTRAP_SERVERS", "localhost:9092")
	schemaRegistryUrl := GetEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081")
	topicName := GetEnv("TOPIC_NAME", "event.customer.entity")
	flag.Parse()
	log.Printf("keySchemaFile: %v\n", *keySchemaFile)
	log.Printf("valueSchemaFile: %v\n", *valueSchemaFile)
	log.Printf("bootstrapServers: %v\n", bootstrapServers)
	log.Printf("schemaRegistryUrl: %v\n", schemaRegistryUrl)
	log.Printf("topicName: %v\n", topicName)

	log.Printf("readSchema(%v)\n", *keySchemaFile)
	keySchema := ReadFile(*keySchemaFile)
	log.Printf("readSchema(%v)\n", *valueSchemaFile)
	valueSchema := ReadFile(*valueSchemaFile)

	log.Printf("BuildProducer\n")
	kafkaProducer, err := BuildProducer(bootstrapServers, schemaRegistryUrl, topicName, keySchema, valueSchema)
	if err != nil {
		panic(err)
	}

	repo := CustomerRepo{}

	for {
		// Dummy: ms-customer -> produce to kafka
		customer, isNew := repo.NextCustomer()
		if isNew {
			kafkaProducer.ProduceCustomer(customer)
		}
		time.Sleep(time.Second * 1)
		// Dummy: core banking -> produce to postgresql
		account := repo.NextAccount(customer)
		time.Sleep(time.Second * 1)
		repo.CreateBooking(account)
		time.Sleep(time.Second * 1)

		// Dummy: Contact Center -> produce to MySql
		c := repo.NextCase(customer)
		time.Sleep(time.Second * 1)
		repo.CreateCommunication(c)
	}
}
