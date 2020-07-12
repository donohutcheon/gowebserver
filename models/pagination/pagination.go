package pagination

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
)

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type OptionalInt64 struct {
	Value int64
	Valid bool
}

func (o *OptionalInt64) Set(value int64) {
	o.Value = value
	o.Valid = true
}

type OptionalString struct {
	Value string
	Valid bool
}

func (o *OptionalString) Set(value string) {
	o.Value = value
	o.Valid = true
}

type SortDirection2 OptionalString



type Parameters struct {
	FetchFrom  OptionalInt64
	Page       OptionalInt64
	FetchCount OptionalInt64
	SortField  string
	SortDir    SortDirection
	isInfinite bool
}

type Sortable interface {
	GetSortFields() map[string]bool
	GetPagination() Parameters
	SetSortParameters(Parameters)
}

func ParsePagination(logger *log.Logger, queryParams url.Values, entity Sortable) error {
	var page OptionalInt64
	var fetchCount OptionalInt64
	var fetchFrom OptionalInt64
	sortField := "id"
	sortDir := SortDirectionAsc
	isInfinite := true

	if _, ok := queryParams["from"]; ok {
		value, err := strconv.ParseInt(queryParams.Get("from"), 10, 64)
		if err != nil {
			return err
		}
		if value < 0 {
			fields := []types.ErrorField{
				{
					Name:    "from",
					Message: "negative from id value",
					Direct:  true,
				},
			}
			return errors.NewError("invalid index parameters", fields, http.StatusBadRequest )
		}
		fetchFrom.Set(value)
	}

	if _, ok := queryParams["page"]; ok {
		value, err := strconv.ParseInt(queryParams.Get("page"), 10, 64)
		if err != nil {
			return err
		}
		if value < 0 {
			fields := []types.ErrorField{
				{
					Name:    "page",
					Message: "negative page value",
					Direct:  true,
				},
			}
			return errors.NewError("invalid pagination parameters", fields, http.StatusBadRequest )
		}
		isInfinite = false
		page.Set(value)
	}

	if _, ok := queryParams["count"]; ok {
		value, err := strconv.ParseInt(queryParams.Get("count"), 10, 64)
		if err != nil {
			return err
		}
		if value < 0 {
			fields := []types.ErrorField{
				{
					Name:    "count",
					Message: "invalid page count value",
					Direct:  true,
				},
			}
			return errors.NewError("invalid pagination parameters", fields, http.StatusBadRequest )
		}
		fetchCount.Set(value)
	}

	if _, ok := queryParams["sortField"]; ok {
		sortField = queryParams.Get("sortField")
		_, ok := entity.GetSortFields()[sortField]
		if !ok {
			fields := []types.ErrorField{
				{
					Name:    "sortField",
					Message: "invalid sort field",
					Direct:  true,
				},
			}
			return errors.NewError("invalid sort field", fields, http.StatusBadRequest )
		}
	}

	if _, ok := queryParams["sortDir"]; ok {
		sortDir = SortDirection(queryParams.Get("sortDir"))
		if sortDir != SortDirectionAsc && sortDir != SortDirectionDesc {
			fields := []types.ErrorField{
				{
					Name:    "sortDir",
					Message: "invalid sort direction",
					Direct:  true,
				},
			}
			return errors.NewError("invalid sort direction", fields, http.StatusBadRequest )
		}
	}

	entity.SetSortParameters(
		Parameters{
			FetchFrom:  fetchFrom,
			Page:       page,
			FetchCount: fetchCount,
			SortField:  sortField,
			SortDir:    sortDir,
			isInfinite: isInfinite,
		})

	return nil
}

func (p *Parameters) BuildPagination(sortColumn string) string {
	var offset int64
	copy := *p
	if !copy.Page.Valid {
		copy.Page.Value = 0
	}
	if !p.FetchCount.Valid {
		copy.FetchCount.Value = 10
	}
	offset = copy.Page.Value * copy.FetchCount.Value
	return fmt.Sprintf(" order by %s %s, id %s limit %d, %d", sortColumn, p.SortDir, p.SortDir, offset, copy.FetchCount.Value)
}