package collector

import "dedicatus545-be/core"

type File interface {
	FileInfo

	// Checks that file is a directory
	IsDir() bool

	// List of files in current directory if IsDir return true, in other case return empty list
	List() []File
}

type FileInfo interface {
	// File name
	Name() string

	// URL to the File that can be used to download file
	Path() string
}

func CollectFromTo(from File, to core.Repository) error {
	return nil
}
