package mockdatalayer

import (
	"database/sql"
	"github.com/donohutcheon/gowebserver/models/filters"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models/pagination"
)

func (m *MockDataLayer) getNextCardTransactionID() int64 {
	var maxID int64 = -1
	for _, cardTransaction := range m.CardTransactions {
		if cardTransaction.ID > maxID {
			maxID = cardTransaction.ID
		}
	}

	return maxID + 1
}

func (m *MockDataLayer) CreateCardTransaction(cardTransaction *datalayer.CardTransaction) (int64, error) {
	cardTransaction.CreatedAt = datalayer.JsonNullTime{
		NullTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	cardTransaction.ID = m.getNextCardTransactionID()

	m.CardTransactions = append(m.CardTransactions, cardTransaction)

	return cardTransaction.ID, nil
}

func (m *MockDataLayer) GetCardTransactionByID(id int64) (*datalayer.CardTransaction, error) {
	for _, cardTransaction := range m.CardTransactions {
		if id == cardTransaction.ID {
			return cardTransaction, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) GetCardTransactionsByUserID(userID int64, sortable pagination.Sortable, filter filters.CardTransactionFilter) ([]*datalayer.CardTransaction, error) {
	var cardTransactions []*datalayer.CardTransaction
	var cardTransaction *datalayer.CardTransaction
	for _, cardTransaction = range m.CardTransactions {
		if userID == cardTransaction.UserID {
			cardTransactions = append(cardTransactions, cardTransaction)
		}
	}

	if len(cardTransactions) == 0 {
		return nil, datalayer.ErrNoData
	}

	return cardTransactions, nil
}