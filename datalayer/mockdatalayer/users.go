package mockdatalayer

import (
	"database/sql"
	"math"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)


func (m *MockDataLayer) GetUserByEmail(email string) (*datalayer.User, error) {
	for _, user := range m.Users {
		if user.Email.Valid && email == user.Email.String {
			return user, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) GetUserByID(id int64) (*datalayer.User, error) {
	for _, user := range m.Users {
		if id == user.ID {
			return user, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) getNextUserID() int64 {
	var maxID int64 = math.MinInt64
	for _, user := range m.Users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	return maxID + 1
}

func (m *MockDataLayer) CreateUser(email, password string) (int64, error){
	user, err := m.GetUserByEmail(email)
	if err != datalayer.ErrNoData {
		return 0, err
	}

	user = &datalayer.User{
		Model:    datalayer.Model{
			ID:        m.getNextUserID(),
			CreatedAt: datalayer.JsonNullTime{
				NullTime: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			},
			UpdatedAt: datalayer.JsonNullTime{},
			DeletedAt: datalayer.JsonNullTime{},
		},
		Email:    sql.NullString{
			String: email,
			Valid:  true,
		},
		Password: sql.NullString{
			String: password,
			Valid:  true,
		},
		State: sql.NullString{
			String: "CONFIRMED",
			Valid:  true,
		},
	}

	m.Users = append(m.Users, user)

	return user.ID, nil
}

// TODO: Implement!
func (m *MockDataLayer) GetUnconfirmedUsers() ([]datalayer.User, error) {
	return nil, nil
}

func (m *MockDataLayer) SetUserStateByID(int64, datalayer.UserState) error {
 return nil
}