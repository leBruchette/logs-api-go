package reader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

const DefaultChunkSize = 1 << 10

type LineReader struct {
	File           *os.File
	FileInfo       os.FileInfo
	FilePath       string
	LinesRequested int
	SearchText     string
	Context        *LineReaderContext
}

// LineReaderContext struct to hold the state of our reading/writing
type LineReaderContext struct {
	HasWritten    bool
	WriteCount    int
	Position      int64
	Leftover      string
	Chunk         []byte
	LinesAsBytes  [][]byte
	LinesReversed []string
}

// NewReader constructor that accepts path and query parameters
func NewReader(filePath string, linesRequested int, searchText string) *LineReader {
	return &LineReader{
		FilePath:       filePath,
		LinesRequested: linesRequested,
		SearchText:     searchText,
	}
}

// InitializeReader initialize the reader: obtain a handle on the file, read the files statistics, and set the initial reader context
func (lr *LineReader) InitializeReader() error {
	var err error
	lr.File, err = os.Open(lr.FilePath)
	if err != nil {
		// assume file doesn't exist if error
		return NewFileNotFoundError(lr.FilePath, err)
	}

	lr.FileInfo, err = lr.File.Stat()
	if err != nil {
		return NewFileStatError(lr.FilePath, err)
	}

	lr.Context = &LineReaderContext{
		Position:     lr.FileInfo.Size(),
		Leftover:     "",
		HasWritten:   false,
		Chunk:        make([]byte, DefaultChunkSize),
		LinesAsBytes: make([][]byte, 0),
	}

	return nil
}

// StreamFileDetails stream the file details (name, size, etc) as well as the beginning of the lines array
func (lr *LineReader) StreamFileDetails(c *gin.Context, flusher http.Flusher) error {
	fileInfoJson, err := json.Marshal(lr.fileInfoJson())
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Writer, "{ \"file\": %s, \"lines\": [\n", fileInfoJson)
	flusher.Flush()
	return nil
}

// StreamFileLines stream the lines from the file in reverse until we have the requested number of lines (or read the entire file)
func (lr *LineReader) StreamFileLines(c *gin.Context, flusher http.Flusher) error {
	for lr.isLinesRemaining() {
		fmt.Printf("Reading chunk from position %d...\n", lr.Context.Position)
		lines, err := lr.readLinesFromChunk()
		if err != nil {
			return err
		}
		for _, line := range lines {
			if lr.isLinePrintable(line) {
				encodedLine, _ := json.Marshal(line)
				// if/else to ensure we don't print a comma before the first line in our output json
				if !lr.Context.HasWritten {
					fmt.Fprintf(c.Writer, "%s", encodedLine)
					lr.Context.HasWritten = true
				} else {
					fmt.Fprintf(c.Writer, ",\n%s", encodedLine)
				}
				flusher.Flush()
				lr.Context.WriteCount++
				if lr.Context.WriteCount >= lr.LinesRequested {
					break
				}
			}
		}
	}
	return nil
}

// CloseFile close the file handle
func (lr *LineReader) CloseFile() func() {
	return func() {
		lr.File.Close()
	}
}

// readLinesFromChunk read the lines from the chunk of the file into a byte array.
// once a chunk is read, split the chunk by newline char and reverse the order while castint to strings
// while executing, update the Context struct for any possible future reads
func (lr *LineReader) readLinesFromChunk() ([]string, error) {
	// inline basic min function
	var readSize = func(a, b int64) int64 {
		if a < b {
			return a
		}
		return b
	}(int64(DefaultChunkSize), lr.Context.Position)
	if readSize < 0 {
		return nil, NewFileSeekError(lr.FilePath, errors.New("Chunk read size is negative"))
	}

	lr.Context.Position -= readSize

	_, err := lr.File.Seek(lr.Context.Position, io.SeekStart)
	if err != nil {
		return nil, NewFileSeekError(lr.FilePath, err)
	}

	lr.Context.Chunk = lr.Context.Chunk[:readSize]
	_, err = lr.File.Read(lr.Context.Chunk)
	if err != nil {
		return nil, NewFileReadError(lr.FilePath, err)
	}

	lr.Context.LinesAsBytes = lr.Context.LinesAsBytes[:0]
	// filter out empty array elements before appending
	for _, line := range bytes.Split(append(lr.Context.Chunk, []byte(lr.Context.Leftover)...), []byte("\n")) {
		if len(line) > 0 {
			lr.Context.LinesAsBytes = append(lr.Context.LinesAsBytes, line)
		}
	}

	if lr.Context.Position > 0 {
		lr.Context.Leftover = string(lr.Context.LinesAsBytes[0])
		lr.Context.LinesAsBytes = lr.Context.LinesAsBytes[1:]
	}

	reversedLines := make([]string, len(lr.Context.LinesAsBytes))
	for i, line := range lr.Context.LinesAsBytes {
		reversedLines[len(lr.Context.LinesAsBytes)-1-i] = string(line)
	}
	return reversedLines, nil
}

// fileInfoJson convenience function for printing "file" json attribute
func (lr *LineReader) fileInfoJson() map[string]interface{} {
	return map[string]interface{}{
		"name":    lr.FilePath,
		"size":    lr.FileInfo.Size(),
		"mode":    lr.FileInfo.Mode(),
		"modTime": lr.FileInfo.ModTime(),
	}
}

// isLinePrintable only print non-empty lines, and lines that contain the search text if provided
func (lr *LineReader) isLinePrintable(line string) bool {
	return len(line) > 0 && strings.TrimSpace(line) != "" && (lr.SearchText == "" || strings.Contains(strings.ToLower(line), strings.ToLower(lr.SearchText)))
}

// isLinesRemaining check if there are lines remaining to be read based on read position and requested line count
func (lr *LineReader) isLinesRemaining() bool {
	return lr.Context.Position > 0 && lr.Context.WriteCount <= lr.LinesRequested
}
