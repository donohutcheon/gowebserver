package datalayer

import (
	"database/sql"
	"fmt"
	"log"
)

type Contact struct {
	Model
	Name   string `json:"name" db:"name"`
	Phone  string `json:"phone" db:"phone"`
	UserID int64  `json:"userID" db:"user_id"`
}

func (p *PersistenceDataLayer) CreateContact(name, phone string, userID int64) (int64, error) {
	result, err := p.GetConn().Exec("insert into contacts(name, phone, user_id) values (?, ?, ?)", name, phone, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return 0, err
	}

	return id, nil
}

func (p *PersistenceDataLayer) GetContactByID(id int64) (*Contact, error) {
	contact := new(Contact)

	row := p.GetConn().QueryRow(`SELECT id, user_id, name, phone, created_at, updated_at, deleted_at FROM contacts WHERE id=?`, id)
	err := row.Scan(contact)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return contact, nil
}

func (p *PersistenceDataLayer) GetContactsByUserID(userID int64) ([]*Contact, error) {
	contacts := make([]*Contact, 0)

	rows, err := p.GetConn().Query(`SELECT id, name, phone, created_at, updated_at, deleted_at FROM contacts WHERE user_id=?`, userID)
	if err == sql.ErrNoRows {
		fmt.Println(false, "User account does not exist. Please re-login")
		return nil, ErrNoData
	} else if err != nil {
		fmt.Printf("Failed to query contacts for user ID [%d] from database", userID)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		contact := new(Contact)
		rows.Scan(contact)
		contacts = append(contacts, contact)
	}

	return contacts, nil
}