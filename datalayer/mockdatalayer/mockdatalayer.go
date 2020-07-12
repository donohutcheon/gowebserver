package mockdatalayer

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

type MockDataLayer struct {
	t                   *testing.T
	Users               []*datalayer.User
	Contacts            []*datalayer.Contact
	CardTransactions    []*datalayer.CardTransaction
	SignUpConfirmations []*datalayer.SignUpConfirmation
	usersFilename       string
	contactsFilename    string
	cardTransFilename   string
	signUpConfsFilename string
}

func New(t *testing.T) *MockDataLayer {
	m := new(MockDataLayer)
	m.t = t
	m.initialize()
	return m
}

func (m *MockDataLayer) initialize() error {
	err := m.LoadUserTestData("testdata/users.json")
	require.NoError(m.t, err)

	err = m.LoadContactTestData("testdata/contacts.json")
	require.NoError(m.t, err)

	return nil
}

func (m *MockDataLayer) ResetAndReload() error {
	m.Users = m.Users[:0]
	err := m.LoadUserTestData(m.usersFilename)
	if err != nil {
		return err
	}

	m.Contacts = m.Contacts[:0]
	err = m.LoadContactTestData(m.contactsFilename)
	if err != nil {
		return err
	}

	m.CardTransactions = m.CardTransactions[:0]

	return nil
}

func (m *MockDataLayer) LoadUserTestData(filename string) error{
	m.usersFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.Users)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadContactTestData(filename string) error{
	m.contactsFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.Contacts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadCardTransactionTestData(filename string) error{
	m.cardTransFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.CardTransactions)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadSignUpConfTestData(filename string) error{
	m.signUpConfsFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.SignUpConfirmations)
	if err != nil {
		return err
	}

	return nil
}