package mockdatalayer

import (
	"github.com/donohutcheon/gowebserver/datalayer"
	"math"
)

func (m *MockDataLayer) CreateSignUpConfirmation(nonce string, userID int64) (int64, error) {
	id := m.getNextSignUpConfID()
	signUpConf := datalayer.SignUpConfirmation{
		Model:  datalayer.Model{
			ID:        id,
			CreatedAt: datalayer.JsonNullTime{},
			UpdatedAt: datalayer.JsonNullTime{},
			DeletedAt: datalayer.JsonNullTime{},
		},
		Nonce:  nonce,
		UserID: userID,
	}
	m.SignUpConfirmations = append(m.SignUpConfirmations, &signUpConf)

	return id, nil
}

func (m *MockDataLayer) LookupSignUpConfirmation(nonce string) (*datalayer.SignUpConfirmation, error) {
	for _, s := range m.SignUpConfirmations {
		if nonce == s.Nonce {
			return s, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) getNextSignUpConfID() int64 {
	var maxID int64 = math.MinInt64
	for _, s := range m.SignUpConfirmations {
		if s.ID > maxID {
			maxID = s.ID
		}
	}

	return maxID + 1
}