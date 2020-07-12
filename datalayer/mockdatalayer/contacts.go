package mockdatalayer

import (
	"database/sql"
	"math"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)

func (m *MockDataLayer) getNextContactID() int64 {
	var maxID int64 = math.MinInt64
	for _, contact := range m.Contacts {
		if contact.ID > maxID {
			maxID = contact.ID
		}
	}

	return maxID + 1
}

func (m *MockDataLayer) CreateContact(name, phone string, userID int64) (int64, error) {
	contact := &datalayer.Contact{
		Model:    datalayer.Model{
			ID:        m.getNextContactID(),
			CreatedAt: datalayer.JsonNullTime{
				NullTime: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			},
			UpdatedAt: datalayer.JsonNullTime{},
			DeletedAt: datalayer.JsonNullTime{},
		},
		Name: name,
		Phone: phone,
	}

	m.Contacts = append(m.Contacts, contact)

	return contact.ID, nil
}

func (m *MockDataLayer) GetContactByID(id int64) (*datalayer.Contact, error) {
	for _, contact := range m.Contacts {
		if id == contact.ID {
			return contact, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) GetContactsByUserID(userID int64) ([]*datalayer.Contact, error) {
	var contacts []*datalayer.Contact
	var contact *datalayer.Contact
	for _, contact = range m.Contacts {
		if userID == contact.UserID {
			contacts = append(contacts, contact)
		}
	}

	if len(contacts) == 0 {
		return nil, datalayer.ErrNoData
	}

	return contacts, nil
}