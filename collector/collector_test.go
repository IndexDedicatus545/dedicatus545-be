package collector

import (
	"errors"
	"testing"
)

func TestCollectFromTo(t *testing.T) {
	tests := []struct {
		name    string
		from    File
		to      Repository
		wantErr bool
	}{
		// TODO: Add test cases.
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
	savedFiles []FileInfo
	returnErr  bool
}

func (t *TestRepository) Save(info FileInfo) error {
	t.savedFiles = append(t.savedFiles, info)
	if t.returnErr {
		return errors.New("failed to save file")
	}
	return nil
}

type TestFile struct {
	title string
	ext   string
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

func (t TestFile) Title() string {
	return t.title
}

func (t TestFile) Extension() string {
	return t.ext
}

func (t TestFile) Path() string {
	return t.path
}

func SingleFile(name string, ext string) File {
	return TestFile{title: name, ext: ext}
}

func Dir(name string, files ...File) File {
	return TestFile{title: name, isDir: true, files: files}
}
