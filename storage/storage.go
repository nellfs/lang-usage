package main

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccoutn(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) erro
}
