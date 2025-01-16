package reader

import "fmt"

type FileError struct {
	FilePath string
	Err      error
}

type FileNotFoundError struct {
	fileError *FileError
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found: %s, error: %v", e.fileError.FilePath, e.fileError.Err)
}

func NewFileNotFoundError(filePath string, err error) *FileNotFoundError {
	return &FileNotFoundError{
		fileError: &FileError{
			FilePath: filePath,
			Err:      err,
		},
	}
}

type FileStatError struct {
	fileError *FileError
}

func (e *FileStatError) Error() string {
	return fmt.Sprintf("failed to stat file: %s, error: %v", e.fileError.FilePath, e.fileError.Err)
}

func NewFileStatError(filePath string, err error) *FileStatError {
	return &FileStatError{
		fileError: &FileError{
			FilePath: filePath,
			Err:      err,
		},
	}
}

type FileSeekError struct {
	fileError *FileError
}

func (e *FileSeekError) Error() string {
	return fmt.Sprintf("failed to seek file: %s, error: %v", e.fileError.FilePath, e.fileError.Err)
}

func NewFileSeekError(filePath string, err error) *FileSeekError {
	return &FileSeekError{
		fileError: &FileError{
			FilePath: filePath,
			Err:      err,
		},
	}
}

type FileReadError struct {
	fileError *FileError
}

func (e *FileReadError) Error() string {
	return fmt.Sprintf("failed to read file: %s, error: %v", e.fileError.FilePath, e.fileError.Err)
}

func NewFileReadError(filePath string, err error) *FileReadError {
	return &FileReadError{
		fileError: &FileError{
			FilePath: filePath,
			Err:      err,
		},
	}
}
