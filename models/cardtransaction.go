package models

import (
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"github.com/donohutcheon/gowebserver/models/filters"
	"github.com/donohutcheon/gowebserver/models/pagination"
	"github.com/donohutcheon/gowebserver/state"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)

type CurrencyValue struct {
	Value int64 `json:"value" db:"amount"`
	Scale int   `json:"scale" db:"currency_scale"`
}

type CardTransaction struct {
	datalayer.Model
	DateTime             time.Time     `json:"dateTime" db:"datetime"`
	Amount               CurrencyValue `json:"amount"`
	CurrencyCode         string        `json:"currencyCode" db:"currency_code"`
	Reference            string        `json:"reference" db:"reference"`
	MerchantName         string        `json:"merchantName" db:"merchant_name"`
	MerchantCity         string        `json:"merchantCity" db:"merchant_city"`
	MerchantCountryCode  string        `json:"merchantCountryCode" db:"merchant_country_code"`
	MerchantCountryName  string        `json:"merchantCountryName" db:"merchant_country_name"`
	MerchantCategoryCode string        `json:"merchantCategoryCode" db:"merchant_category_code"`
	MerchantCategoryName string        `json:"merchantCategoryName" db:"merchant_category_name"`
	UserID               int64         `json:"userID" db:"user_id"`
	serverState          *state.ServerState
	pagination           pagination.Parameters
	filter               filters.CardTransactionFilter
}



func (c *CardTransaction) GetSortFields() map[string]bool {
	return map[string]bool {
		"id": true,
		"amount" : true,
		"currencyCodes" : true,
		"dateTime" : true,
		"references" : true,
		"merchantNames" : true,
		"merchantCities" : true,
		"merchantCountryCodes" : true,
		"merchantCountryNames" : true,
		"merchantCategoryCodes" : true,
		"merchantCategoryNames" : true,
	}
}

func (c *CardTransaction) SetSortParameters(parameters pagination.Parameters) {
	c.pagination = parameters
}

func (c *CardTransaction) GetPagination() pagination.Parameters {
	return c.pagination
}

func NewCardTransaction(state *state.ServerState) *CardTransaction {
	cardTransaction := new(CardTransaction)
	cardTransaction.serverState = state

	return cardTransaction
}

func newFromDBCardTransaction(cardTransaction *datalayer.CardTransaction) *CardTransaction{
	c := new(CardTransaction)
	c.ID = cardTransaction.ID
	c.CreatedAt = cardTransaction.CreatedAt
	c.UpdatedAt = cardTransaction.UpdatedAt
	c.DeletedAt = cardTransaction.DeletedAt
	c.DateTime = cardTransaction.DateTime
	c.Amount.Value = cardTransaction.Amount
	c.Amount.Scale = cardTransaction.CurrencyScale
	c.CurrencyCode = cardTransaction.CurrencyCode
	c.Reference = cardTransaction.Reference
	c.MerchantName = cardTransaction.MerchantName
	c.MerchantCity = cardTransaction.MerchantCity
	c.MerchantCountryCode = cardTransaction.MerchantCountryCode
	c.MerchantCountryName = cardTransaction.MerchantCountryName
	c.MerchantCategoryCode = cardTransaction.MerchantCategoryCode
	c.MerchantCategoryName = cardTransaction.MerchantCategoryName
	return c
}

func (c *CardTransaction) convertToDB() *datalayer.CardTransaction {
	cardTransaction := new(datalayer.CardTransaction)
	cardTransaction.ID = c.ID
	cardTransaction.CreatedAt = c.CreatedAt
	cardTransaction.UpdatedAt = c.UpdatedAt
	cardTransaction.DeletedAt = c.DeletedAt
	cardTransaction.DateTime = c.DateTime
	cardTransaction.Amount = c.Amount.Value
	cardTransaction.CurrencyScale = c.Amount.Scale
	cardTransaction.CurrencyCode = c.CurrencyCode
	cardTransaction.Reference = c.Reference
	cardTransaction.MerchantName = c.MerchantName
	cardTransaction.MerchantCity = c.MerchantCity
	cardTransaction.MerchantCountryCode = c.MerchantCountryCode
	cardTransaction.MerchantCountryName = c.MerchantCountryName
	cardTransaction.MerchantCategoryCode = c.MerchantCategoryCode
	cardTransaction.MerchantCategoryName = c.MerchantCategoryName
	cardTransaction.UserID = c.UserID
	return cardTransaction
}

// TODO: return errors
func (c *CardTransaction) validate() error {
	// TODO: Validate no empty fields

	if c.UserID <= 0 {
		return ErrUserDoesNotExist
	}

	if len(c.CurrencyCode) == 0 {
		return ErrValidationFailed
	}

	if len(c.MerchantName) == 0 {
		return ErrValidationFailed
	}

	//All the required parameters are present
	return nil
}

func (c *CardTransaction) CreateCardTransaction() (*CardTransaction, error) {
	// TODO: c.Validate() to return an error
	err := c.validate()
	if err != nil {
		return nil, err
	}

	dbCardTransaction := c.convertToDB()

	dl := c.serverState.DataLayer
	id, err := dl.CreateCardTransaction(dbCardTransaction)
	if err != nil {
		c.serverState.Logger.Println(err)
		return nil, err
	}

	dbCardTransaction, err = dl.GetCardTransactionByID(id)
	if err != nil {
		return nil, err
	}

	data := newFromDBCardTransaction(dbCardTransaction)

	return data, nil
}

func (c *CardTransaction) GetCardTransaction(id int64) (*CardTransaction, error) {
	dl := c.serverState.DataLayer
	dbCardTransaction, err := dl.GetCardTransactionByID(id)
	if err == datalayer.ErrNoData {
		return nil, err // TODO: return proper error with code
	}

	cardTransaction := newFromDBCardTransaction(dbCardTransaction)

	return cardTransaction, nil
}

func (c *CardTransaction) GetCardTransactionsByUserID(userID int64) ([]*CardTransaction, error) {
	dl := c.serverState.DataLayer
	cardTransactions := make([]*CardTransaction, 0)

	dbCardTransactions, err := dl.GetCardTransactionsByUserID(userID, c, c.filter)
	if err == datalayer.ErrNoData {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	for _, dbCardTransaction := range dbCardTransactions {
		cardTransaction := newFromDBCardTransaction(dbCardTransaction)
		cardTransactions = append(cardTransactions, cardTransaction)
	}

	return cardTransactions, err
}

func (c *CardTransaction) SetFilterCriteria(queryParams url.Values) error {
	err := c.filterAmount(queryParams)
	if err != nil {
		return err
	}

	err = c.filterDateTime(queryParams)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardTransaction) filterAmount(queryParams url.Values) error {
	queryVal, ok := queryParams["amount"]
	if !ok {
		return nil
	}

	if len(queryVal) == 0 {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid value"},
		}, http.StatusBadRequest)
	}
	amountRange := strings.Split(queryVal[0], "-")
	if len(amountRange) != 2 {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid amount range"},
		}, http.StatusBadRequest)
	}

	value, err := strconv.ParseInt(amountRange[0], 10, 64)
	if err != nil {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid lowerbound amount value"},
		}, http.StatusBadRequest)
	}
	c.filter.Amount.LowerBound = value

	value, err = strconv.ParseInt(amountRange[1], 10, 64)
	if err != nil {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid upperbound amount value"},
		}, http.StatusBadRequest)
	}
	c.filter.Amount.UpperBound = value

	if c.filter.Amount.UpperBound < c.filter.Amount.LowerBound {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "upperbound must be greater than the lowerbound amount value"},
		}, http.StatusBadRequest)
	}
	c.filter.Amount.IsSet = true

	return nil
}

func (c *CardTransaction) filterDateTime(queryParams url.Values) error {
	queryVal, ok := queryParams["dateTime"]
	if !ok {
		return nil
	}

	if len(queryVal) == 0 {
		return e.NewError("dateTime filter range is invalid", []types.ErrorField{
			{Name: "dateTime", Message: "invalid value"},
		}, http.StatusBadRequest)
	}

	dateRange := strings.Split(queryVal[0], "-")
	if len(dateRange) != 2 {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid amount range"},
		}, http.StatusBadRequest)
	}

	lowerValue, err := strconv.ParseInt(dateRange[0], 10, 64)
	if err != nil {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid lowerbound amount value"},
		}, http.StatusBadRequest)
	}
	c.filter.DateTime.LowerBound = time.Unix(lowerValue, 0)

	upperValue, err := strconv.ParseInt(dateRange[1], 10, 64)
	if err != nil {
		return e.NewError("amount filter range is invalid", []types.ErrorField{
			{Name: "amount", Message: "invalid upperbound amount value"},
		}, http.StatusBadRequest)
	}
	c.filter.DateTime.UpperBound = time.Unix(upperValue, 0)

	if upperValue < lowerValue {
		return e.NewError("dateTime filter range is invalid", []types.ErrorField{
			{Name: "dateTime", Message: "upperbound must be greater than the lowerbound dateTime value"},
		}, http.StatusBadRequest)
	}
	c.filter.DateTime.IsSet = true

	return nil
}