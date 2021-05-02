package collector

import (
	"dedicatus545-be/core"
	"fmt"
	"github.com/alecthomas/units"
	"strconv"
	"strings"
)

type File interface {
	FileInfo

	// Checks that file is a directory
	IsDir() (bool, error)

	// List of files in current directory if IsDir return true, in other case return empty list
	List() ([]File, error)
}

type FileInfo interface {
	// File name
	Name() string

	// File size in bytes
	Size() uint64

	// URL to the File that can be used to download file
	Path() string
}

//todo сделать красивые ошибки
func CollectFromTo(from File, to core.Repository) error {
	files, err := flatAllFiles(from)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}

	books := make([]*core.Book, 0, len(files))
	for _, file := range files {
		book, err := createBookFrom(file)
		if err != nil {
			return fmt.Errorf("failed to create book: %w", err)
		}
		books = append(books, book)
	}

	for _, book := range books {
		if err := to.Save(*book); err != nil {
			return fmt.Errorf("failed to save collected file: %s: %w", strconv.Quote(book.Title), err)
		}
	}

	return nil
}

func flatAllFiles(from File) ([]FileInfo, error) {
	isDir, err := from.IsDir()
	if err != nil {
		return nil, fmt.Errorf("failed on %s: %w", strconv.Quote(from.Name()), err)
	}
	if !isDir {
		return []FileInfo{from}, nil
	}

	list, err := from.List()
	if err != nil {
		return nil, fmt.Errorf("failed on %s: %w", strconv.Quote(from.Name()), err)
	}
	var dirs []File
	dirs = append(dirs, list...)

	var result []FileInfo
	i := 0
	for {
		if i >= len(dirs) {
			break
		}
		isDir, err := dirs[i].IsDir()
		if err != nil {
			return nil, fmt.Errorf("failed on %s: %w", strconv.Quote(dirs[i].Name()), err)
		}
		if isDir {
			list, err := dirs[i].List()
			if err != nil {
				return nil, fmt.Errorf("failed on %s: %w", strconv.Quote(dirs[i].Name()), err)
			}
			dirs = append(dirs, list...)
		} else {
			result = append(result, dirs[i])
		}

		i++
	}

	return result, nil
}

func createBookFrom(file FileInfo) (*core.Book, error) {
	index := strings.LastIndex(file.Name(), ".")
	var name string
	var extension string
	if index == -1 {
		name = file.Name()
	} else {
		name = file.Name()[:index]
		extension = file.Name()[index+1:]
	}

	bookFile, err := core.NewBookFile(file.Path(), extension, units.MetricBytes(file.Size()))
	if err != nil {
		return nil, fmt.Errorf("failed to create book file: %s: %w", strconv.Quote(name), err)
	}
	book, err := core.NewBook(name, bookFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create book from file: %s: %w", strconv.Quote(name), err)
	}
	return book, nil
}
