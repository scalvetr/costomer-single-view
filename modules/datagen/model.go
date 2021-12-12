package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CustomerStruct struct {
	CustomerId string            `json:"customerId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Surname    string            `json:"surname,omitempty"`
	Email      string            `json:"email,omitempty"`
	Telephones []TelephoneStruct `json:"telephones,omitempty"`
	Addresses  []AddressStruct   `json:"addresses,omitempty"`
}

type AddressStruct struct {
	Street  string `json:"street,omitempty"`
	Number  string `json:"number,omitempty"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	ZipCode string `json:"zipCode,omitempty"`
	Default bool   `json:"default"`
}
type TelephoneStruct struct {
	Number  string `json:"number,omitempty"`
	Primary bool   `json:"primary"`
}

type AccountStatus int64

const (
	Open AccountStatus = iota
	Cancelled
)

func (s AccountStatus) String() string {
	switch s {
	case Open:
		return "OPEN"
	case Cancelled:
		return "CANCELLED"
	}
	return "unknown"
}

type AccountStruct struct {
	AccountId        int32         `json:"accountId,omitempty" db:"account_id"`
	CustomerId       string        `json:"customerId,omitempty" db:"customer_id"`
	IBAN             string        `json:"iban,omitempty" db:"iban"`
	Balance          float64       `json:"balance,omitempty" db:"balance"`
	CreationDate     time.Time     `json:"creationDate,omitempty" db:"creation_date"`
	CancellationDate *time.Time    `json:"cancellationDate,omitempty" db:"cancellation_date"`
	Status           AccountStatus `json:"status,omitempty" db:"status"`
}

type BookingStruct struct {
	BookingId   int32     `json:"bookingId,omitempty"`
	AccountId   int32     `json:"accountId,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Description string    `json:"description,omitempty"`
	BookingDate time.Time `json:"bookingDate,omitempty"`
	ValueDate   time.Time `json:"valueDate,omitempty"`
	Fee         float64   `json:"fee,omitempty"`
	Taxes       float64   `json:"taxes,omitempty"`
}

type CaseStruct struct {
	ID                primitive.ObjectID        `bson:"_id,omitempty"`
	CaseId            string                    `bson:"case_id,omitempty"`
	CustomerId        string                    `bson:"customer_id,omitempty"`
	Title             string                    `bson:"title,omitempty"`
	CreationTimestamp primitive.DateTime        `bson:"creation_timestamp,omitempty"`
	Communications    []CaseCommunicationStruct `bson:"communications,omitempty"`
}
type CaseCommunicationStruct struct {
	CommunicationId string             `bson:"communication_id,omitempty"`
	Type            string             `bson:"type,omitempty"`
	Text            string             `bson:"text,omitempty"`
	Notes           string             `bson:"notes,omitempty"`
	Timestamp       primitive.DateTime `bson:"timestamp,omitempty"`
}
