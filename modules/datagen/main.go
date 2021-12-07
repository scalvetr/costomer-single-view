package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	keySchemaFile := flag.String("key-schema-file", "customer-key.avsc", "AVRO key schema file")
	valueSchemaFile := flag.String("value-schema-file", "customer-value.avsc", "AVRO value schema file")
	// kafka
	bootstrapServers := GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	schemaRegistryUrl := GetEnv("KAFKA_SCHEMA_REGISTRY_URL", "http://localhost:8081")
	topicName := GetEnv("KAFKA_TOPIC_NAME", "event.customer.entity")
	// core banking
	coreBankingDbHost := GetEnv("CORE_BANKING_DB_HOST", "localhost")
	coreBankingDbPort := GetEnv("CORE_BANKING_DB_PORT", "5432")
	coreBankingDbUser := GetEnv("CORE_BANKING_DB_USER", "user")
	coreBankingDbPassword := GetEnv("CORE_BANKING_DB_PASSWORD", "password")
	coreBankingDbName := GetEnv("CORE_BANKING_DB_NAME", "core-banking")
	// core banking
	contactCenterDbUri := GetEnv("CONTACT_CENTER_DB_URI", "mongodb://localhost:27017")
	contactCenterDbUser := GetEnv("CONTACT_CENTER_DB_USER", "user")
	contactCenterDbPassword := GetEnv("CONTACT_CENTER_DB_PASSWORD", "password")
	contactCenterDbName := GetEnv("CONTACT_CENTER_DB_NAME", "core-banking")

	flag.Parse()
	log.Printf("keySchemaFile: %v\n", *keySchemaFile)
	log.Printf("valueSchemaFile: %v\n", *valueSchemaFile)
	log.Printf("bootstrapServers: %v\n", bootstrapServers)
	log.Printf("schemaRegistryUrl: %v\n", schemaRegistryUrl)
	log.Printf("topicName: %v\n", topicName)
	log.Printf("coreBankingDbHost: %v\n", coreBankingDbHost)
	log.Printf("coreBankingDbPort: %v\n", coreBankingDbPort)
	log.Printf("coreBankingDbUser: %v\n", coreBankingDbUser)
	log.Printf("coreBankingDbPassword: %v\n", coreBankingDbPassword)
	log.Printf("coreBankingDbName: %v\n", coreBankingDbName)
	log.Printf("contactCenterDbUri: %v\n", contactCenterDbUri)
	log.Printf("contactCenterDbUser: %v\n", contactCenterDbUser)
	log.Printf("contactCenterDbPassword: %v\n", contactCenterDbPassword)
	log.Printf("contactCenterDbName: %v\n", contactCenterDbName)

	log.Printf("readSchema(%v)\n", *keySchemaFile)
	keySchema := ReadFile(*keySchemaFile)
	log.Printf("readSchema(%v)\n", *valueSchemaFile)
	valueSchema := ReadFile(*valueSchemaFile)

	log.Printf("BuildProducer\n")
	kafkaProducer, err := BuildProducer(bootstrapServers, schemaRegistryUrl, topicName, keySchema, valueSchema)
	if err != nil {
		panic(err)
	}

	repo := BuildCustomerRepo(PgDbConfig{
		DbHost:     coreBankingDbHost,
		DbPort:     coreBankingDbPort,
		DbUser:     coreBankingDbUser,
		DbPassword: coreBankingDbPassword,
		DbName:     coreBankingDbName,
	}, MongoDbConfig{
		DbUri:      contactCenterDbUri,
		DbUser:     contactCenterDbUser,
		DbPassword: contactCenterDbPassword,
		DbName:     contactCenterDbName,
	})

	for {
		// Dummy: ms-customer -> produce to kafka
		customer, isNew := repo.NextCustomer()
		log.Printf("NextCustomer() = %v, %v\n", customer, isNew)
		if isNew {
			kafkaProducer.ProduceCustomer(customer)
			log.Printf("ProduceCustomer(%v)\n", customer)
		}
		time.Sleep(time.Second * 1)
		// Dummy: core banking -> produce to postgresql
		account := repo.NextAccount(customer)
		log.Printf("NextAccount(%v) = %v\n", customer, account)
		time.Sleep(time.Second * 1)
		booking := repo.CreateBooking(account)
		log.Printf("CreateBooking(%v) = %v\n", account, booking)
		time.Sleep(time.Second * 1)

		// Dummy: Contact Center -> produce to MySql
		c := repo.NextCase(customer)
		time.Sleep(time.Second * 1)
		repo.CreateCommunication(c)
	}
}
