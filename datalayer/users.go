package datalayer

import (
	"database/sql"
)

type User struct {
	Model
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
	Role     sql.NullString `db:"role"`
	State     sql.NullString `db:"state"`
	LoggedOutAt JsonNullTime `db:"logged_out_at"`
}

func (p *PersistenceDataLayer) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := p.GetConn().QueryRowx(`select * from users where email = ?`, email)
	err := row.StructScan(user)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PersistenceDataLayer) GetUserByID(id int64) (*User, error) {
	user := new(User)
	row := p.GetConn().QueryRowx(`SELECT * FROM users WHERE id=?`, id)
	err := row.StructScan(user)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PersistenceDataLayer) CreateUser(email, password string) (int64, error){
	result, err := p.GetConn().Exec("insert into users(email, password, state) values (?, ?, ?)", email, password, UserStateUnconfirmed)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PersistenceDataLayer) GetUnconfirmedUsers() ([]User, error) {
	var users []User
	err := p.GetConn().Select(&users, `SELECT * FROM users WHERE state=?`, UserStateUnconfirmed)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PersistenceDataLayer) SetUserStateByID(id int64, state UserState) error {
	result, err := p.GetConn().Exec("update users set state = ? where id = ?", state, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrNoData
	}

	return nil
}