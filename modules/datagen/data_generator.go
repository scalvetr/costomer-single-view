package main

import (
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"strings"
	"time"
)

var length = 0

const capacity = 1000

var customers = new([capacity]CustomerStruct)

type DataGenerator struct {
	coreBankingRepo   CoreBankingRepo
	contactCenterRepo ContactCenterRepo
	random            *rand.Rand
	faker             faker.Faker
}

func BuildDataGenerator(dbConfig PgDbConfig, mongoConfig MongoDbConfig) DataGenerator {
	return DataGenerator{
		coreBankingRepo:   BuildCoreBankingRepo(dbConfig),
		contactCenterRepo: BuildContactCenterRepo(mongoConfig),
		random:            rand.New(rand.NewSource(22)),
		faker:             faker.New(),
	}
}

func (g DataGenerator) generateCustomer() CustomerStruct {
	faker := faker.New()
	firstName := faker.Person().FirstName()
	lastName := faker.Person().LastName()
	address := faker.Address()

	data := CustomerStruct{
		CustomerId: strings.ToLower(firstName + "." + lastName),
		Email:      strings.ToLower(firstName + "." + lastName + "@gmail.com"),
		Name:       firstName,
		Surname:    lastName,
		Telephones: []TelephoneStruct{
			{Number: faker.Phone().Number(), Primary: true},
			{Number: faker.Phone().Number(), Primary: false},
		},
		Addresses: []AddressStruct{
			{Number: address.BuildingNumber(),
				City:    address.City(),
				Country: address.Country(),
				Street:  address.StreetName(),
				ZipCode: address.PostCode(),
				Default: true,
			},
		},
	}
	return data
}

func (g DataGenerator) NextCustomer() (CustomerStruct, bool) {
	// until the array reaches its capacity
	if length < capacity {
		data := g.generateCustomer()
		customers[length] = data
		length++
		return data, true
	} else {
		createNew := rand.Int()%2 == 0
		index := rand.Intn(length - 1)
		if createNew {
			customers[index] = g.generateCustomer()
		}
		return customers[index], true
	}
}

func (g DataGenerator) NextAccount(customer CustomerStruct) AccountStruct {
	createNew := rand.Int()%2 == 0
	if !createNew {
		account := g.coreBankingRepo.GetOpenAccount(customer.CustomerId)
		if account != nil {
			log.Printf("return an existing account with id (%v)\n", account.AccountId)
			return *account
		}
	}
	log.Printf("Creating a new one for customer %v\n", customer.CustomerId)
	faker := faker.New()
	account := g.coreBankingRepo.StoreAccount(AccountStruct{
		CustomerId:       customer.CustomerId,
		AccountId:        0,
		Balance:          0,
		IBAN:             faker.Payment().CreditCardNumber(),
		CreationDate:     time.Now(),
		CancellationDate: nil,
		Status:           Open,
	})
	return account
}

func (g DataGenerator) NextBooking(account AccountStruct) BookingStruct {
	accountId := account.AccountId
	amount := g.random.Float64() * 10.000

	booking := g.coreBankingRepo.StoreBooking(BookingStruct{
		BookingId:   0,
		AccountId:   accountId,
		Amount:      amount,
		Description: g.faker.Lorem().Sentence(g.random.Intn(15)),
		BookingDate: time.Now(),
		ValueDate:   time.Now(),
		Fee:         0.0,
		Taxes:       0.0,
	})
	g.coreBankingRepo.UpdateAccountBalance(accountId, account.Balance+booking.Amount)
	return booking
}

func (g DataGenerator) NextCase(customer CustomerStruct) CaseStruct {
	createNew := rand.Int()%2 == 0
	if !createNew {
		c := g.contactCenterRepo.GetOpenCase(customer.CustomerId)
		if c != nil {
			log.Printf("return an existing account with id (%v)\n", customer.CustomerId)
			return *c
		}
	}
	c := CaseStruct{
		CaseId:            uuid.New().String(),
		CustomerId:        customer.CustomerId,
		Title:             g.faker.Lorem().Sentence(g.random.Intn(15)),
		CreationTimestamp: primitive.NewDateTimeFromTime(time.Now()),
		Communications:    []CaseCommunicationStruct{},
	}
	return c
}

func (g DataGenerator) NextCommunication(c CaseStruct) (primitive.ObjectID, CaseCommunicationStruct) {
	communication := CaseCommunicationStruct{
		CommunicationId: uuid.New().String(),
		Text:            g.faker.Lorem().Sentence(g.random.Intn(15)),
		Type:            g.faker.RandomStringElement([]string{"Mobile", "Web", "Phone"}),
		Notes:           g.faker.Lorem().Sentence(g.random.Intn(25)),
		Timestamp:       primitive.NewDateTimeFromTime(time.Now()),
	}
	c.Communications = append(c.Communications, communication)
	caseId := g.contactCenterRepo.StoreCase(c)
	return caseId, communication
}

func (g DataGenerator) Close() {
	err := g.contactCenterRepo.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = g.coreBankingRepo.Close()
	if err != nil {
		log.Fatal(err)
	}
}
