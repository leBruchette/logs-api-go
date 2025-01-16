package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"logs-api-go/reader"
	"net/http"
	"strconv"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	router := configureRouter()
	// Start the server
	router.Run(":8080")
}

func configureRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, &Message{Text: "OK"})
	})

	router.GET("/logs/*filepath", func(c *gin.Context) {
		// Setting headers to indicate streaming
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// Ensure we can flush the response incrementally flusher,
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		// Create a LineReader instance
		lr := createLineReaderForRequest(c)

		// Initialize the reader
		err := lr.InitializeReader()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer lr.CloseFile()()

		// stream the response, starting with file details and opening of our lines array
		err = lr.StreamFileDetails(c, flusher)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// stream the lines from the file in reverse until we have the requested number of lines (or read the entire file)
		err = lr.StreamFileLines(c, flusher)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// close the json
		fmt.Fprintf(c.Writer, "]\n}")
		flusher.Flush()
	})

	return router
}

func createLineReaderForRequest(c *gin.Context) *reader.LineReader {
	linesRequested, _ := strconv.Atoi(c.DefaultQuery("lines", "30"))
	searchText := c.DefaultQuery("search", "")
	lr := reader.NewReader(c.Param("filepath"), linesRequested, searchText)
	return lr
}
