package collector

import (
	"dedicatus545-be/core"
	"errors"
	"github.com/alecthomas/units"
	"testing"
)

func TestCollectFromTo(t *testing.T) {
	tests := []struct {
		name     string
		from     File
		to       core.Repository
		wantErr  bool
		expected []core.Book
	}{
		{
			name:     "collect and save single file",
			from:     TFile("file-1.pdf"),
			to:       MockRepository(),
			expected: []core.Book{
				//core.NewBook("file-1", core.NewBookFile()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CollectFromTo(tt.from, tt.to); (err != nil) != tt.wantErr {
				t.Errorf("CollectFromTo() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

type TestRepository struct {
	savedBooks []core.Book
	returnErr  bool
}

func (t *TestRepository) Save(book core.Book) error {
	t.savedBooks = append(t.savedBooks, book)
	if t.returnErr {
		return errors.New("failed to save file")
	}
	return nil
}

func MockRepository() core.Repository {
	return &TestRepository{}
}

type TestFile struct {
	name  string
	isDir bool
	files []File
	path  string
}

func (t TestFile) IsDir() bool {
	return t.isDir
}

func (t TestFile) List() []File {
	return t.files
}

func (t TestFile) Name() string {
	return t.name
}

func (t TestFile) Path() string {
	return t.path
}

func TFile(name string) File {
	return TestFile{name: name}
}

func Dir(name string, files ...File) File {
	return TestFile{name: name, isDir: true, files: files}
}

func Book(name, ext, path string, size units.SI) core.Book {
	file, err := core.NewBookFile(path, ext, size)
	if err != nil {
		panic(ext)
	}
	book, err := core.NewBook(name, file)
	if err != nil {
		panic(ext)
	}
	return *book
}
