package main

import (
	"github.com/jaswdr/faker"
	"log"
	"math/rand"
	"strings"
	"time"
)

var length = 0

const capacity = 1000

var customers = new([capacity]CustomerStruct)

type CustomerRepo struct {
	PgDbConfig      PgDbConfig
	coreBankingRepo CoreBankingRepo
}

func (r CustomerRepo) buildCustomer() CustomerStruct {
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
	//var data CustomerStruct
	//err := faker.FakeData(&data)
	//if err != nil {
	//	panic(err)
	//}
	return data
}

func (r CustomerRepo) NextCustomer() (CustomerStruct, bool) {
	if length < capacity {
		data := r.buildCustomer()
		customers[length] = data
		length++
		return data, true
	} else {
		createNew := rand.Int()%2 == 0
		index := rand.Intn(length - 1)
		if createNew {
			customers[index] = r.buildCustomer()
		}
		return customers[index], true
	}
}
func BuildCustomerRepo(dbConfig PgDbConfig) CustomerRepo {
	return CustomerRepo{
		PgDbConfig:      dbConfig,
		coreBankingRepo: BuildCoreBankingRepo(dbConfig),
	}
}

func (r CustomerRepo) NextAccount(customer CustomerStruct) AccountStruct {
	createNew := rand.Int()%2 == 0
	if !createNew {
		account := r.coreBankingRepo.GetOpenAccount(customer.CustomerId)
		if account != nil {
			log.Printf("return an existing account with id (%v)\n", account.AccountId)
			return *account
		}
	}
	faker := faker.New()
	account := r.coreBankingRepo.StoreAccount(AccountStruct{
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

func (r CustomerRepo) CreateBooking(account AccountStruct) BookingStruct {
	return BookingStruct{}
}

func (r CustomerRepo) NextCase(customer CustomerStruct) CaseStruct {
	return CaseStruct{}
}

func (r CustomerRepo) CreateCommunication(c CaseStruct) CaseCommunicationStruct {
	return CaseCommunicationStruct{}
}
