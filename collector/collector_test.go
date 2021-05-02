package collector

import (
	"dedicatus545-be/core"
	"errors"
	"github.com/alecthomas/units"
	"reflect"
	"testing"
)

func TestCollectFromTo(t *testing.T) {
	tests := []struct {
		name    string
		from    File
		to      core.Repository
		wantErr bool
		want    []core.Book
	}{
		{
			name: "collect and save single file",
			from: TestFile{name: "file-1.pdf", path: "test://file-1.pdf"},
			to:   MockRepository(),
			want: []core.Book{
				MustBook("file-1", "pdf", "test://file-1.pdf", 0),
			},
		},
		{
			name: "collect file with size",
			from: TestFile{name: "file-1.pdf", path: "test://file-1.pdf", size: 14000000},
			to:   MockRepository(),
			want: []core.Book{
				MustBook("file-1", "pdf", "test://file-1.pdf", 14*units.MB),
			},
		},
		{
			name: "collect file without extension",
			from: TestFile{name: "file-1", path: "test://file-1"},
			to:   MockRepository(),
			want: []core.Book{
				MustBook("file-1", "", "test://file-1", 0),
			},
		},
		{
			name: "collect and save one file in directory",
			from: Dir("foo", TestFile{name: "file-1.pdf", path: "test://file-1.pdf"}),
			to:   MockRepository(),
			want: []core.Book{
				MustBook("file-1", "pdf", "test://file-1.pdf", 0),
			},
		},
		{
			name: "not collect directory as book",
			from: Dir("foo"),
			to:   MockRepository(),
		},
		{
			name: "collect several files in directory",
			from: Dir("foo",
				TestFile{name: "file-1.pdf", path: "test://file-1.pdf"},
				TestFile{name: "file-2.pdf", path: "test://file-2.pdf"}),
			to: MockRepository(),
			want: []core.Book{
				MustBook("file-1", "pdf", "test://file-1.pdf", 0),
				MustBook("file-2", "pdf", "test://file-2.pdf", 0),
			},
		},
		//todo собрать несколько файлов как одну книгу
		//todo не сохранять книги которые уже есть (проверка уникальности)
		//todo не прекращать сборку, если пару книг оказались испорченными
		//todo собирать в несколько потоков
		{
			name: "collect from 2 deep dir",
			from: Dir("foo",
				Dir("bar",
					TestFile{name: "file-1.pdf", path: "test://file-1.pdf"},
					TestFile{name: "file-2.pdf", path: "test://file-2.pdf"})),
			to: MockRepository(),
			want: []core.Book{
				MustBook("file-1", "pdf", "test://file-1.pdf", 0),
				MustBook("file-2", "pdf", "test://file-2.pdf", 0),
			},
		},
		{
			name: "collect files from all levels",
			from: Dir("foo",
				TestFile{name: "file-1.pdf", path: "test://file-1.pdf"},
				Dir("bar",
					Dir("bar-2",
						TestFile{name: "file-2.pdf", path: "test://file-2.pdf"}),
					TestFile{name: "file-3.pdf", path: "test://file-3.pdf"}),
				Dir("baz",
					TestFile{name: "file-4.pdf", path: "test://file-4.pdf"}),
			),
			to: MockRepository(),
			want: []core.Book{
				// в порядке обхода дерева в ширину
				MustBook("file-1", "pdf", "test://file-1.pdf", 0),
				MustBook("file-3", "pdf", "test://file-3.pdf", 0),
				MustBook("file-4", "pdf", "test://file-4.pdf", 0),
				MustBook("file-2", "pdf", "test://file-2.pdf", 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CollectFromTo(tt.from, tt.to); (err != nil) != tt.wantErr {
				t.Fatalf("CollectFromTo() error = %v, wantErr %v", err, tt.wantErr)
			}

			r := tt.to.(*TestRepository)
			if !reflect.DeepEqual(r.savedBooks, tt.want) {
				t.Errorf("CollectFromTo() got = %v, want %v", r.savedBooks, tt.want)
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
	size  uint64

	returnErr bool
}

func (t TestFile) Size() uint64 {
	return t.size
}

func (t TestFile) IsDir() (bool, error) {
	if t.returnErr {
		return t.isDir, errors.New("IsDir() failed")
	} else {
		return t.isDir, nil
	}
}

func (t TestFile) List() ([]File, error) {
	if t.returnErr {
		return t.files, errors.New("List() failed")
	} else {
		return t.files, nil
	}
}

func (t TestFile) Name() string {
	return t.name
}

func (t TestFile) Path() string {
	return t.path
}

func Dir(name string, files ...File) File {
	return TestFile{name: name, isDir: true, files: files}
}

func MustBook(name, ext, path string, size units.MetricBytes) core.Book {
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
