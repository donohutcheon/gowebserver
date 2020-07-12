package models

type Sortable interface {
	GetSortFields() map[string]bool
}
