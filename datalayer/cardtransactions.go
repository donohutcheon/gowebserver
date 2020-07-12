package datalayer

import (
	"database/sql"
	"fmt"
	"github.com/donohutcheon/gowebserver/models/filters"
	"github.com/donohutcheon/gowebserver/models/pagination"
	"log"
	"strings"
	"time"
)

type CardTransaction struct {
	Model
	DateTime             time.Time `json:"dateTime" db:"datetime"`
	Amount               int64     `json:"amount" db:"amount"`
	CurrencyScale        int       `json:"scale" db:"currency_scale"`
	CurrencyCode         string    `json:"currencyCode" db:"currency_code"`
	Reference            string    `json:"reference" db:"reference"`
	MerchantName         string    `json:"merchantName" db:"merchant_name"`
	MerchantCity         string    `json:"merchantCity" db:"merchant_city"`
	MerchantCountryCode  string    `json:"merchantCountryCode" db:"merchant_country_code"`
	MerchantCountryName  string    `json:"merchantCountryName" db:"merchant_country_name"`
	MerchantCategoryCode string    `json:"merchantCategoryCode" db:"merchant_category_code"`
	MerchantCategoryName string    `json:"merchantCategoryName" db:"merchant_category_name"`
	UserID               int64     `json:"userID" db:"user_id"`
}



func (p *PersistenceDataLayer) CreateCardTransaction(cardTransaction *CardTransaction) (int64, error) {
	const cols = "datetime, amount, currency_scale, currency_code, reference, merchant_name, merchant_city, merchant_country_code, merchant_country_name, merchant_category_code, merchant_category_name, user_id"
	var bindCols = ":" + strings.ReplaceAll(cols, ", ", ", :")

	sql := fmt.Sprintf("insert into card_transactions(%s) values (%s)", cols, bindCols)
	result, err := p.GetConn().NamedExec(sql, cardTransaction)
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

func (p *PersistenceDataLayer) GetCardTransactionByID(id int64) (*CardTransaction, error) {
	cardTransaction := new(CardTransaction)
	statement := "SELECT * FROM card_transactions WHERE id=?"
	row := p.GetConn().QueryRowx(statement, id)
	err := row.StructScan(cardTransaction)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return cardTransaction, nil
}

func (p *PersistenceDataLayer) GetCardTransactionsByUserID(userID int64, sortable pagination.Sortable, filter filters.CardTransactionFilter) ([]*CardTransaction, error) {
	cardTransactions := make([]*CardTransaction, 0)
	pageParams := sortable.GetPagination()
	filterSQL, filterValues := GetFilterCriteria(filter)

	dbSortField := pageParams.SortField

	/*var pageFilterDir string
	if pageParams.SortDir == pagination.SortDirectionDesc {
		pageFilterDir = "<"
	} else {
		pageFilterDir = ">"
	}*/

	//pagination := fmt.Sprintf(" and id %s ? order by %s %s, id %s limit %d", pageFilterDir, dbSortField, pageParams.SortDir, pageParams.SortDir, pageParams.FetchCount)
	//offset := pageParams.Page * pageParams.FetchCount
	pagination := pageParams.BuildPagination(dbSortField)
	//pagination := fmt.Sprintf(" order by %s %s, id %s limit %d, %d", dbSortField, pageParams.SortDir, pageParams.SortDir, offset, pageParams.FetchCount)
	statement := "SELECT * FROM card_transactions WHERE user_id=? " + filterSQL + pagination
	fmt.Println(statement)
	var bindValues []interface{}
	bindValues = append(bindValues, userID)
	bindValues = append(bindValues, filterValues...)
	//bindValues = append(bindValues, pageParams.FetchFrom)
	rows, err := p.GetConn().Queryx(statement, bindValues...)
	if err == sql.ErrNoRows {
		fmt.Printf("There are no card transactions for user ID [%d]", userID)
		return nil, ErrNoData
	} else if err != nil {
		fmt.Printf("Failed to query card transactions for user ID [%d] from database", userID)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cardTransaction := new(CardTransaction)
		rows.StructScan(cardTransaction)
		cardTransactions = append(cardTransactions, cardTransaction)
	}

	return cardTransactions, nil
}

func GetFilterCriteria(filter filters.CardTransactionFilter) (string, []interface{}) {
	builder := new(strings.Builder)
	var values []interface{}
	if filter.Amount.IsSet {
		predicate := " and amount >= ? and amount < ? "
		builder.WriteString(predicate)
		values = append(values, filter.Amount.LowerBound)
		values = append(values, filter.Amount.UpperBound)
	}

	if filter.DateTime.IsSet {
		predicate := " and datetime >= ? and datetime < ? "
		builder.WriteString(predicate)
		values = append(values, filter.DateTime.LowerBound)
		values = append(values, filter.DateTime.UpperBound)
	}

	return builder.String(), values
}