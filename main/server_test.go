package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T) *os.File {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	content := "line1\nline2\nline3\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return tempFile
}

func setup(t *testing.T) (*gin.Engine, *os.File) {
	router := configureRouter()
	tempFile := createTestFile(t)
	return router, tempFile
}

func teardown(tempFile *os.File) {
	os.Remove(tempFile.Name())
}

func TestLogsEndpointNoParams(t *testing.T) {
	router, tempFile := setup(t)
	defer teardown(tempFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s", tempFile.Name()), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), fmt.Sprintf("\"name\":\"%s\"", tempFile.Name()))
}

func TestLogsEndpointInvalidFilePath(t *testing.T) {
	router := configureRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs/invalidfile.log?lines=5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "file not found")
}

func TestLogsEndpointWithLinesParam(t *testing.T) {
	router, tempFile := setup(t)
	defer teardown(tempFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s?lines=2", tempFile.Name()), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"lines\": [\n\"line3\",\n\"line2\"]\n")
}

func TestLogsEndpointWithSearchText(t *testing.T) {
	router, tempFile := setup(t)
	defer teardown(tempFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s?lines=5&search=line1", tempFile.Name()), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"lines\": [\n\"line1\"]\n")
}
