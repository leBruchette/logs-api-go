package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"logs-api-go/reader"
	"net/http"
	"strconv"
)

type Message struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

func main() {
	router := configureRouter()
	// Start the server
	router.Run(":8080")
}

func configureRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, &Message{Error: "OK"})
	})

	router.GET("/logs/*filepath", func(c *gin.Context) {
		// Setting headers to indicate streaming
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// Ensure we can flush the response incrementally flusher,
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, &Message{Error: "Streaming unsupported!"})
			return
		}

		// Create a LineReader instance
		lr := createLineReaderForRequest(c)

		// Initialize the reader
		err := lr.InitializeReader()
		if err != nil {
			handleError(c, err)
			return
		}
		defer lr.CloseFile()()

		// stream the response, starting with file details and opening of our lines array
		err = lr.StreamFileDetails(c, flusher)
		if err != nil {
			handleError(c, err)
			return
		}

		// stream the lines from the file in reverse until we have the requested number of lines (or read the entire file)
		err = lr.StreamFileLines(c, flusher)
		if err != nil {
			handleError(c, err)
		}

		// close the json
		fmt.Fprintf(c.Writer, "]\n}")
		flusher.Flush()
	})

	return router
}

func createLineReaderForRequest(c *gin.Context) *reader.LineReader {
	linesRequested, err := strconv.Atoi(c.DefaultQuery("lines", "30"))
	if err != nil {
		// default to 30 lines if the query parameter is invalid.  Alpha numeric values are converted to zero, so
		// reinforce the default value here.
		linesRequested = 30
	}
	searchText := c.DefaultQuery("search", "")
	lr := reader.NewReader(c.Param("filepath"), linesRequested, searchText)
	return lr
}

func handleError(c *gin.Context, err error) {
	if errors.As(err, &reader.FileNotFoundError) {
		c.JSON(http.StatusNotFound, &Message{Error: err.Error()})
	} else if errors.As(err, &reader.FileStatError) {
		c.JSON(http.StatusInternalServerError, &Message{Error: err.Error()})
	} else if errors.As(err, &reader.FileSeekError) {
		c.JSON(http.StatusInternalServerError, &Message{Error: err.Error()})
	} else if errors.As(err, &reader.FileReadError) {
		c.JSON(http.StatusInternalServerError, &Message{Error: err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, &Message{Error: err.Error()})
	}
	return
}
