package core

type Repository interface {
	Save(book Book) error
}
