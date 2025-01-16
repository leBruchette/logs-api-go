package reader

import (
	"os"
	"testing"
)

func TestInitializeReaderSuccess(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	content := "line1\nline2\nline3\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	file.Close()

	lr := NewReader(file.Name(), 5, "")
	err = lr.InitializeReader()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestInitializeReaderFileNotFound(t *testing.T) {
	lr := NewReader("nonexistentfile", 5, "")
	err := lr.InitializeReader()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestReadLinesFromChunkSuccess(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	content := "line1\nline2\nline3\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	file.Close()

	lr := NewReader(file.Name(), 5, "")
	err = lr.InitializeReader()
	if err != nil {
		t.Fatalf("failed to initialize reader: %v", err)
	}

	lines, err := lr.readLinesFromChunk()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	expectedLines := []string{"line3", "line2", "line1"}
	for i, line := range expectedLines {
		if lines[i] != line {
			t.Errorf("expected line: %s, got: %s", line, lines[i])
		}
	}
}

func TestReadLinesFromChunkFileSeekError(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	content := "line1\nline2\nline3\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	file.Close()

	lr := NewReader(file.Name(), 5, "")
	err = lr.InitializeReader()
	if err != nil {
		t.Fatalf("failed to initialize reader: %v", err)
	}

	lr.Context.Position = -1
	_, err = lr.readLinesFromChunk()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIsLinePrintableSuccess(t *testing.T) {
	lr := NewReader("", 5, "line2")
	if !lr.isLinePrintable("line2") {
		t.Errorf("expected line to be printable")
	}
}

func TestIsLinePrintableEmptyLine(t *testing.T) {
	lr := NewReader("", 5, "")
	if lr.isLinePrintable("") {
		t.Errorf("expected line to be not printable")
	}
}
