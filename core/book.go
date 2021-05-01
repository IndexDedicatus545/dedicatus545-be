package core

import (
	"github.com/alecthomas/units"
	valid "github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

//todo считаю, что никто не будет создавать напрямую через структуру,
// а не через специальную функцию, но вообще надо бы сделать защиту от этого
type Book struct {
	ID        string
	ISBN      ISBN
	Title     string
	Desc      string
	PageCount uint32

	Authors []Author
	Files   []BookFile
	Tags    []Tag

	valid bool
}

type ISBN string

type BookFile struct {
	Path   string
	Format string
	Size   units.MetricBytes

	valid bool
}

type Author struct {
	ID  string
	FIO string
}

type Tag string

func NewBook(title string, file BookFile) (*Book, error) {
	if strings.TrimSpace(title) == "" {
		return nil, newInvalidBookError("title must not be blank")
	}
	if !file.valid {
		return nil, newInvalidBookError("book file must be valid")
	}

	b := &Book{Title: title, valid: true}
	b.Files = append(b.Files, file)
	return b, nil
}

func NewBookFile(path, format string, size units.MetricBytes) (BookFile, error) {
	if strings.TrimSpace(path) == "" {
		return BookFile{}, newInvalidBookError("path to book must not be blank")
	}
	if u, err := url.Parse(path); err != nil || !u.IsAbs() {
		return BookFile{}, newInvalidBookErrorWithClause("path to book must be in absolute URL format: ["+path+"]", err)
	}
	return BookFile{
		Path:   path,
		Format: format,
		Size:   size,
		valid:  true,
	}, nil
}

func NewISBN(raw string) (ISBN, error) {
	if valid.IsISBN10(raw) || valid.IsISBN13(raw) {
		return ISBN(raw), nil
	}
	return "", newInvalidISBNError(raw)
}
