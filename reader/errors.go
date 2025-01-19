package reader

import "fmt"

type FileNotFound struct {
	FilePath string
	Err      error
}

func (e *FileNotFound) Error() string {
	return fmt.Sprintf("file not found: %s, error: %v", e.FilePath, e.Err.Error())
}

func NewFileNotFoundError(filePath string, err error) *FileNotFound {

	return &FileNotFound{
		FilePath: filePath,
		Err:      err,
	}
}

var FileNotFoundError = &FileNotFound{}

type FileStat struct {
	FilePath string
	Err      error
}

func (e *FileStat) Error() string {
	return fmt.Sprintf("failed to stat file: %s, error: %v", e.FilePath, e.Err)
}

func NewFileStatError(filePath string, err error) *FileStat {
	return &FileStat{
		FilePath: filePath,
		Err:      err,
	}

}

var FileStatError = &FileStat{}

type FileSeek struct {
	FilePath string
	Err      error
}

func (e *FileSeek) Error() string {
	return fmt.Sprintf("failed to seek file: %s, error: %v", e.FilePath, e.Err)
}

func NewFileSeekError(filePath string, err error) *FileSeek {
	return &FileSeek{
		FilePath: filePath,
		Err:      err,
	}
}

var FileSeekError = &FileSeek{}

type FileRead struct {
	FilePath string
	Err      error
}

func (e *FileRead) Error() string {
	return fmt.Sprintf("failed to read file: %s, error: %v", e.FilePath, e.Err)
}

func NewFileReadError(filePath string, err error) *FileRead {
	return &FileRead{
		FilePath: filePath,
		Err:      err,
	}
}

var FileReadError = &FileRead{}
