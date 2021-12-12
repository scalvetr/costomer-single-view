package main

import (
	"database/sql"
	"fmt"
	"github.com/blockloop/scan"
	_ "github.com/lib/pq"
	"log"
)

type PgDbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}
type CoreBankingRepo struct {
	db       *sql.DB
	dbConfig PgDbConfig
}

func BuildCoreBankingRepo(dbConfig PgDbConfig) CoreBankingRepo {

	info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbName)

	internalDb, err := sql.Open("postgres", info)

	if err != nil {
		panic(err)
	}
	return CoreBankingRepo{
		dbConfig: dbConfig,
		db:       internalDb,
	}
}

func (r CoreBankingRepo) GetOpenAccount(customerId string) *AccountStruct {

	rows, err := r.db.Query(`SELECT 
       account_id, 
       customer_id, 
       iban, balance, 
       creation_date, 
       cancellation_date, 
       status 
	FROM accounts 
	WHERE customer_id = $1 
	AND status = 'OPEN'
	ORDER BY RANDOM() LIMIT 1`, customerId)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var account AccountStruct
		scan.Row(&account, rows)
		return &account
	}
	defer rows.Close()
	return nil
}

func (r CoreBankingRepo) StoreAccount(account AccountStruct) AccountStruct {
	err := r.db.QueryRow(`INSERT 
		INTO accounts (
					  customer_id,
					  iban,
					  balance,
					  creation_date,
					  cancellation_date,
					  status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING account_id`,
		account.CustomerId,
		account.IBAN,
		account.Balance,
		account.CreationDate,
		account.CancellationDate,
		account.Status.String()).Scan(&account.AccountId)
	if err != nil {
		panic(err)
	}
	log.Printf("Created a new account with id=%v\n", account.AccountId)
	return account

}

func (r CoreBankingRepo) GetAccount(accountId int32) AccountStruct {
	rows, err := r.db.Query(`SELECT * FROM accounts
		WHERE account_id = $1`,
		accountId)

	if err != nil {
		panic(err)
	}
	var account AccountStruct
	scan.Row(&account, rows)
	defer rows.Close()
	return account

}
func (r CoreBankingRepo) UpdateAccountBalance(accountId int32, balance float64) AccountStruct {
	rows, err := r.db.Query(`UPDATE accounts 
	    SET balance = $1
	    WHERE account_id = $2`,
		balance,
		accountId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	return r.GetAccount(accountId)

}

func (r CoreBankingRepo) StoreBooking(booking BookingStruct) BookingStruct {
	err := r.db.QueryRow(`INSERT 
			INTO bookings (
						  account_id,
						  amount,
						  description,
						  booking_date,
						  value_date,
						  fee,
						  taxes)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING booking_id`,
		booking.AccountId,
		booking.Amount,
		booking.Description,
		booking.BookingDate,
		booking.ValueDate,
		booking.Fee,
		booking.Taxes,
	).Scan(&booking.BookingId)
	if err != nil {
		panic(err)
	}
	return booking
}

func (r CoreBankingRepo) Close() error {
	err := r.db.Close()
	if err == nil {
		fmt.Println("Connection to Postgresql closed.")
	}
	return err
}
