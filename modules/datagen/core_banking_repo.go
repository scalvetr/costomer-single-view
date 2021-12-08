package main

import (
	"database/sql"
	"fmt"
	"github.com/blockloop/scan"
	_ "github.com/lib/pq"
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
	FROM account 
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
	//defer rows.Close()
	return nil
}

func (r CoreBankingRepo) StoreAccount(account AccountStruct) AccountStruct {
	r.db.QueryRow(`INSERT 
		INTO account (
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
	return account

}

func (r CoreBankingRepo) GetAccount(accountId int32) AccountStruct {
	rows, err := r.db.Query(`SELECT * FROM account
		WHERE account_id = $1`,
		accountId)

	if err != nil {
		panic(err)
	}
	var account AccountStruct
	scan.Row(&account, rows)
	//defer rows.Close()
	return account

}
func (r CoreBankingRepo) UpdateAccountBalance(accountId int32, balance float64) AccountStruct {
	_, err := r.db.Query(`UPDATE account 
	    SET balance = $1
	    WHERE account_id = $2`,
		balance,
		accountId)
	if err != nil {
		panic(err)
	}

	return r.GetAccount(accountId)

}

func (r CoreBankingRepo) StoreBooking(booking BookingStruct) BookingStruct {
	err := r.db.QueryRow(`INSERT 
			INTO booking (
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