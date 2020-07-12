package datalayer

import (
	"github.com/donohutcheon/gowebserver/models/filters"
	"github.com/donohutcheon/gowebserver/models/pagination"
)

type UserState string

const (
	UserStateInvalid     UserState = "INVALID"
	UserStateUnconfirmed UserState = "UNCONFIRMED"
	UserStateProcessing  UserState = "PROCESSING"
	UserStatePending     UserState = "PENDING"
	UserStateConfirmed   UserState = "CONFIRMED"
)

type DataLayer interface {
	// Users
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	CreateUser(email, password string) (int64, error)
	GetUnconfirmedUsers() ([]User, error)
	SetUserStateByID(id int64, state UserState) error

	// Contacts
	CreateContact(name, phone string, userID int64) (int64, error)
	GetContactByID(id int64) (*Contact, error)
	GetContactsByUserID(userID int64) ([]*Contact, error)

	// Transactions
	CreateCardTransaction(*CardTransaction) (int64, error)
	GetCardTransactionByID(id int64) (*CardTransaction, error)
	GetCardTransactionsByUserID(userID int64, sortable pagination.Sortable, filter filters.CardTransactionFilter) ([]*CardTransaction, error)

	// SignUpConfirmations
	CreateSignUpConfirmation(nonce string, userID int64) (int64, error)
	LookupSignUpConfirmation(nonce string) (*SignUpConfirmation, error)
}