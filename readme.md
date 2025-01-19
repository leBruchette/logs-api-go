# Logs API Go

This project provides an API for reading log files in reverse order. It includes a GET endpoint that allows users to retrieve log lines based on specific query parameters.

## Endpoints

## GET `/logs/<path>`

### Path Variables
- `<path>` (string, required): The path to the log file from the root directory.

### Query Parameters
- `lines` (int, optional): The number of lines to read in reverse order. Default is 10.
- `search` (string, optional): Text to search for within the log lines. Only lines containing this text will be returned.

### Example Request
```
GET /logs/var/log/syslog&lines=5&search=error
```

### Example Response
```json
{
  "file": {
    "name": "/var/log/syslog",
    "size": 123456,
    "mode": "0644",
    "modTime": "2023-10-01T12:34:56Z"
  },
  "lines": [
    "error: something went wrong",
    "error: another issue occurred"
  ]
}
```

## Implementation Details

### `reader/reader.go`

This file contains the implementation of the `LineReader` struct and its methods for reading log files in reverse order.

#### `LineReader` Struct

```go
type LineReader struct {
 File           *os.File
 FileInfo       os.FileInfo
 FilePath       string
 LinesRequested int
 SearchText     string
 Context        *LineReaderContext
}
```

- `File`: The file handle for the log file.
- `FileInfo`: Metadata about the log file.
- `FilePath`: The path to the log file.
- `LinesRequested`: The number of lines to read.
- `SearchText`: Text to search for within the log lines.
- `Context`: Holds the state of the reading process.

#### `LineReaderContext` Struct

```go
type LineReaderContext struct {
 HasWritten    bool
 Position      int64
 Leftover      string
 Chunk         []byte
 LinesAsBytes  [][]byte
 LinesReversed []string
}
```

- `HasWritten`: Indicates if any lines have been written.
- `Position`: The current position in the file.
- `Leftover`: Any leftover data from the previous read.
- `Chunk`: The current chunk of data being read.
- `LinesAsBytes`: The lines read as byte slices.
- `LinesReversed`: The lines read in reverse order.

#### Methods

- `NewReader(filePath string, linesRequested int, searchText string) *LineReader`: Constructor for `LineReader`.
- `InitializeReader() error`: Initializes the reader by opening the file and reading its metadata.
- `ReadLinesFromChunk() ([]string, error)`: Reads lines from the file in reverse order.
- `FileInfoJson() map[string]interface{}`: Returns file metadata as a JSON object.
- `IsLinePrintable(line string) bool`: Determines if a line should be printed based on the search text.

### Example Usage

```go
lr := reader.NewReader("/var/log/syslog", 5, "error")
err := lr.InitializeReader()
if err != nil {
    log.Fatalf("failed to initialize reader: %v", err)
}

lines, err := lr.ReadLinesFromChunk()
if err != nil {
    log.Fatalf("failed to read lines: %v", err)
}

for _, line := range lines {
    if lr.IsLinePrintable(line) {
        fmt.Println(line)
    }
}
```

This example initializes a `LineReader`, reads lines from the log file, and prints the most recent 5 lines that contain the search text "error".  The lines are read and processed in chunks, allowing for efficient handling of large files.

```bash

### Make commands

The `Makefile` includes the following targets for starting the application and running tests.

- `verify-deps`: Verifies all dependencies are installed, currently only checks `docker-compose`.
- `all`: The default target that runs the `start` target.
- `start`: Runs the application by executing `main/server.go`.
- `test`: Runs all tests in the project with colored output using `gotest`.
- `docker-build`: Builds a Docker image using `docker-compose`.
- `docker-run`: Runs the Docker container on port 8100 using `docker-compose`.
