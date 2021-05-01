package core

import (
	"github.com/alecthomas/units"
	"reflect"
	"testing"
)

func TestNewBookFile(t *testing.T) {
	type args struct {
		path   string
		format string
		size   units.MetricBytes
	}
	tests := []struct {
		name    string
		args    args
		want    BookFile
		wantErr bool
	}{
		{
			name: "create book file with all fields",
			args: args{
				path:   "test://path/to/book.pdf",
				format: "pdf",
				size:   33 * units.MB,
			},
			want: BookFile{
				Path:   "test://path/to/book.pdf",
				Format: "pdf",
				Size:   33 * units.MB,
			},
		},
		{
			name: "create book file without extension",
			args: args{
				path: "test://path/to/book.pdf",
				size: 33 * units.MB,
			},
			want: BookFile{
				Path: "test://path/to/book.pdf",
				Size: 33 * units.MB,
			},
		},
		{
			name: "create book file without extension and size",
			args: args{
				path: "test://path/to/book.pdf",
			},
			want: BookFile{
				Path: "test://path/to/book.pdf",
			},
		},
		{
			name:    "book file required path",
			args:    args{},
			wantErr: true,
		},
		{
			name: "book file required path in valid URL format",
			args: args{
				path: "hello_world",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookFile(tt.args.path, tt.args.format, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBookFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBook(t *testing.T) {
	type args struct {
		title string
		file  BookFile
	}
	tests := []struct {
		name    string
		args    args
		want    *Book
		wantErr bool
	}{
		{
			name: "create book with all fields",
			args: args{
				title: "book-1",
				file:  MustBookFile("test://to/the/book"),
			},
			want: &Book{
				Title: "book-1",
				Files: []BookFile{
					MustBookFile("test://to/the/book"),
				},
			},
		},
		{
			name: "book required at least one file",
			args: args{
				title: "book-1",
			},
			wantErr: true,
		},
		{
			name: "book required title",
			args: args{
				file: MustBookFile("test://to/the/book"),
			},
			wantErr: true,
		},
		{
			name: "book required not blank title",
			args: args{
				title: "  ",
				file:  MustBookFile("test://to/the/book"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBook(tt.args.title, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func MustBookFile(path string) BookFile {
	file, err := NewBookFile(path, "", 0)
	if err != nil {
		panic(err)
	}
	return file
}
