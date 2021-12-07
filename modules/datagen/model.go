package main

import "time"

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

type AccountStruct struct {
	AccountId        int32     `json:"accountId,omitempty"`
	CustomerId       string    `json:"customerId,omitempty"`
	IBAN             string    `json:"iban,omitempty"`
	Balance          float64   `json:"balance,omitempty"`
	CreationDate     time.Time `json:"creationDate,omitempty"`
	CancellationDate time.Time `json:"cancellationDate,omitempty"`
	Status           string    `json:"status,omitempty"`
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
	CaseId            string    `json:"caseId,omitempty"`
	CustomerId        string    `json:"customerId,omitempty"`
	CreationTimestamp time.Time `json:"creation_timestamp,omitempty"`
}
type CaseCommunicationStruct struct {
	CommunicationId int64  `json:"communicationId,omitempty"`
	CaseId          string `json:"caseId,omitempty"`
	Type            string `json:"type,omitempty"`
	Text            string `json:"text,omitempty"`
	Notes           string `json:"notes,omitempty"`
}
