package models

import (
	"time"
)

type Customer struct {
	CustomerId string      `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Name       string      `json:"name,omitempty" bson:"name,omitempty"`
	Surname    string      `json:"surname,omitempty" bson:"surname,omitempty"`
	Email      string      `json:"email,omitempty" bson:"email,omitempty"`
	Telephones []Telephone `json:"telephones,omitempty" bson:"telephones,omitempty"`
	Addresses  []Address   `json:"addresses,omitempty" bson:"addresses,omitempty"`
	Cases      []Case      `json:"cases,omitempty" bson:"cases,omitempty"`
	Accounts   []Account   `json:"accounts,omitempty" bson:"accounts,omitempty"`
}

type Address struct {
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
	Number  string `json:"number,omitempty" bson:"number,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`
	ZipCode string `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
	Default bool   `json:"default" bson:"default"`
}
type Telephone struct {
	Number  string `json:"number,omitempty" bson:"number,omitempty"`
	Primary bool   `json:"primary" bson:"primary"`
}

type Account struct {
	AccountId        int32      `json:"account_id,omitempty" bson:"account_id,omitempty"`
	CustomerId       string     `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	IBAN             string     `json:"iban,omitempty" bson:"iban,omitempty"`
	Balance          float64    `json:"balance,omitempty" bson:"balance,omitempty"`
	CreationDate     time.Time  `json:"creation_date" bson:"creation_date,omitempty"`
	CancellationDate *time.Time `json:"cancellation_date" bson:"cancellation_date,omitempty"`
	Status           string     `json:"status,omitempty" bson:"status,omitempty"`
}

type Booking struct {
	BookingId   int32     `json:"booking_id,omitempty" bson:"booking_id,omitempty"`
	AccountId   int32     `json:"account_id,omitempty" bson:"account_id,omitempty"`
	Amount      float64   `json:"amount,omitempty" bson:"amount,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	BookingDate time.Time `json:"booking_date" bson:"booking_date,omitempty"`
	ValueDate   time.Time `json:"value_date" bson:"value_date,omitempty"`
	Fee         float64   `json:"fee,omitempty" bson:"fee,omitempty"`
	Taxes       float64   `json:"taxes,omitempty" bson:"taxes,omitempty"`
}

type Case struct {
	CaseId            string              `json:"case_id,omitempty" bson:"case_id,omitempty"`
	CustomerId        string              `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Title             string              `json:"title,omitempty" bson:"title,omitempty"`
	CreationTimestamp time.Time           `json:"creation_timestamp" bson:"creation_timestamp,omitempty"`
	Communications    []CaseCommunication `json:"communications,omitempty" bson:"communications,omitempty"`
}
type CaseCommunication struct {
	CommunicationId string    `json:"communication_id,omitempty" bson:"communication_id,omitempty"`
	Type            string    `json:"type,omitempty" bson:"type,omitempty"`
	Text            string    `json:"text,omitempty" bson:"text,omitempty"`
	Notes           string    `json:"notes,omitempty" bson:"notes,omitempty"`
	Timestamp       time.Time `json:"timestamp" bson:"timestamp,omitempty"`
}
