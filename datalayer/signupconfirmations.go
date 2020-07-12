package datalayer

import (
	"database/sql"
	"log"
)

type SignUpConfirmation struct {
	Model
	Nonce  string `json:"name" db:"nonce"`
	UserID int64  `json:"user_id" db:"user_id"`
}

func (p *PersistenceDataLayer) CreateSignUpConfirmation(nonce string, userID int64) (int64, error) {
	result, err := p.GetConn().Exec("insert into sign_up_confirmations(nonce, user_id) values (?, ?)", nonce, userID)
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


func (p *PersistenceDataLayer) LookupSignUpConfirmation(nonce string) (*SignUpConfirmation, error) {
	signUp := new(SignUpConfirmation)
	statement := "SELECT * FROM sign_up_confirmations WHERE nonce=?"
	row := p.GetConn().QueryRowx(statement, nonce)
	err := row.StructScan(signUp)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return signUp, nil
}