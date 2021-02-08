package repository

type SearchRepository interface {
	Search(interface{}) (interface{}, []error)
}