package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	EditAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}
func NewPostgresStore() (*PostgresStore, error) {

	connStr := "host=localhost port=5432 user=postgres password=ali5101990 dbname=expertdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
	id serial primary key,
	prefix varchar,
	expertName varchar,
	affiliation varchar,
	bh varchar,
	available varchar,
	rating varchar,
	expertRole varchar,
	expertType varchar,
	general_area varchar,
	specialised_Area varchar,
	trained varchar,
	primary_contact varchar,
	secondary_contact varchar,
	email varchar,
	published varchar)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account (prefix,
                     expertName, 
                     affiliation, 
                     bh,
                     available,
                     rating,
                     expertRole,
                     expertType,
                     general_area,
                     specialised_area,
                     trained,
                     primary_contact, 
                     secondary_contact,
                     email,
                     published)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	resp, err := s.db.Query(query,
		acc.Prefix,
		acc.Name,
		acc.Affiliation,
		acc.BH,
		acc.Available,
		acc.Rating,
		acc.Role,
		acc.Type,
		acc.GeneralArea,
		acc.SpecialisedArea,
		acc.Trained,
		acc.PrimaryContact,
		acc.SecondaryContact,
		acc.Email,
		acc.Published)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", resp)
	return nil
}
func (s *PostgresStore) EditAccount(*Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.Prefix,
			&account.Name,
			&account.Affiliation,
			&account.BH,
			&account.Available,
			&account.Rating,
			&account.Role,
			&account.Type,
			&account.GeneralArea,
			&account.SpecialisedArea,
			&account.Trained,
			&account.PrimaryContact,
			&account.SecondaryContact,
			&account.Email,
			&account.Published)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
