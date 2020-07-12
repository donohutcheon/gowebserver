package models

import (
	"github.com/donohutcheon/gowebserver/state"
	"log"

	"github.com/donohutcheon/gowebserver/datalayer"
)

type Contact struct {
	datalayer.Model
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	UserID    int64  `json:"userID"`
	serverState *state.ServerState
}

func NewContact(state *state.ServerState) *Contact {
	contact := new(Contact)
	contact.serverState = state
	return contact
}

func (c *Contact) convert(contact *datalayer.Contact) {
	c.ID = contact.ID
	c.CreatedAt = contact.CreatedAt
	c.UpdatedAt = contact.UpdatedAt
	c.DeletedAt = contact.DeletedAt
	c.Name = contact.Name
	c.Phone = contact.Phone
}

// TODO: return errors
func (c *Contact) validate() error {
	if c.Name == "" {
		return ErrValidationName
	}

	if c.Phone == "" {
		return ErrValidationPhone
	}

	if c.UserID <= 0 {
		return ErrUserDoesNotExist
	}

	//All the required parameters are present
	return nil
}

func (c *Contact) Create() (*Contact, error) {
	// TODO: c.Validate() to return an error
	err := c.validate()
	if err != nil {
		return nil, err
	}

	dl := c.serverState.DataLayer
	id, err := dl.CreateContact(c.Name, c.Phone, c.UserID)
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return nil, err
	}

	dbContact, err := dl.GetContactByID(id)
	if err != nil {
		return nil, err
	}
	c.convert(dbContact)

	return c, nil
}

func (c *Contact) GetContact(id int64) (*Contact, error) {
	dl := c.serverState.DataLayer
	dbContact, err := dl.GetContactByID(id)
	if err == datalayer.ErrNoData {
		return nil, err // TODO: return proper error with code
	}

	contact := new(Contact)
	contact.convert(dbContact)

	return contact, nil
}

func (c *Contact) GetContacts(userID int64) ([]*Contact, error) {
	dl := c.serverState.DataLayer
	contacts := make([]*Contact, 0)

	dbContacts, err := dl.GetContactsByUserID(userID)
	if err == datalayer.ErrNoData {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	for _, dbContact := range dbContacts {
		contact := new(Contact)
		contact.convert(dbContact)
		contacts = append(contacts, contact)
	}

	return contacts, err
}
