package collector

type File interface {
	FileInfo

	// Checks that file is a directory
	IsDir() bool

	// List of files in current directory if IsDir return true, in other case return empty list
	List() []File
}

type FileInfo interface {
	// File name without extension
	Title() string

	// File extension in lower case (e.g. pdf, epub, fb2)
	Extension() string

	// URL to the File path that can be used to download file
	Path() string
}

type Repository interface {
	Save(info FileInfo) error
}

func CollectFromTo(from File, to Repository) error {
	return nil
}
