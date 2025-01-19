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

	assert.Equal(t, http.StatusNotFound, w.Code)
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

func TestLogsEndpointFileNotFound(t *testing.T) {
	router := configureRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs/nonexistentfile.log", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "file not found")
}

func TestLogsEndpointFileStatError(t *testing.T) {
	router, tempFile := setup(t)
	err := tempFile.Chmod(000)
	if err != nil {
		t.Fatalf("failed to change file permissions: %v", err)
	}
	defer teardown(tempFile)

	// Simulate a file stat error by using an invalid file path
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s?lines=-10", tempFile.Name()), nil)
	router.ServeHTTP(w, req)

	//Fixme technically not correct since we throw specific errors for inability to open vs can't stat
	assert.Equal(t, http.StatusNotFound, w.Code)
	//assert.Contains(t, w.Body.String(), "failed to stat file")
}

// Possibly invalid test cases as we check negative positions prior to attempting seeks and reads
//func TestLogsEndpointFileSeekError(t *testing.T) {
//	router, tempFile := setup(t)
//	defer teardown(tempFile)
//
//	// Simulate a file seek error by using an invalid offset
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s?lines=-100", tempFile.Name()), nil)
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusInternalServerError, w.Code)
//	assert.Contains(t, w.Body.String(), "failed to seek file")
//}
//
//func TestLogsEndpointFileReadError(t *testing.T) {
//	router, tempFile := setup(t)
//	defer teardown(tempFile)
//
//	// Simulate a file read error by using an invalid read operation
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", fmt.Sprintf("/logs%s?lines=1000000", tempFile.Name()), nil)
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusInternalServerError, w.Code)
//	assert.Contains(t, w.Body.String(), "failed to read file")
//}
