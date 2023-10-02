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
	name varchar(50),
	affiliation varchar(25),
	primary_contact varchar(50),
	created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account (name, affiliation, primary_contact, created_at)
values ($1, $2, $3, $4)`
	resp, err := s.db.Query(query,
		acc.Name,
		acc.Affiliation,
		acc.PrimaryContact,
		acc.CreatedAt)
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
			&account.Name,
			&account.Affiliation,
			&account.PrimaryContact,
			&account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
