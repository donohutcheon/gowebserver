package filters

import "time"

type AmountRange struct {
	LowerBound int64
	UpperBound int64
	Scale      int
	IsSet      bool
}

type DateRange struct {
	LowerBound time.Time
	UpperBound time.Time
	IsSet      bool
}

type StringFilter struct {
	Value []string
	IsSet bool
}

type CardTransactionFilter struct {
	Amount                AmountRange
	CurrencyCodes         StringFilter
	DateTime              DateRange
	References            StringFilter
	MerchantNames         StringFilter
	MerchantCities        StringFilter
	MerchantCountryCodes  StringFilter
	MerchantCountryNames  StringFilter
	MerchantCategoryCodes StringFilter
	MerchantCategoryNames StringFilter
}