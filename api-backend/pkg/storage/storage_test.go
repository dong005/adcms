package storage

import (
	"strings"
	"testing"
)

func TestLocalStorage_UploadReader(t *testing.T) {
	tmpDir := t.TempDir()
	s := NewLocalStorage(tmpDir, "/uploads")

	reader := strings.NewReader("hello world")
	info, err := s.UploadReader(reader, "hello.txt", "test")
	if err != nil {
		t.Fatalf("UploadReader failed: %v", err)
	}
	if info.Size != 11 {
		t.Errorf("Size = %d, want 11", info.Size)
	}
	if info.URL == "" {
		t.Error("URL should not be empty")
	}
}

func TestLocalStorage_GetURL(t *testing.T) {
	s := NewLocalStorage("/tmp/uploads", "/uploads")
	url := s.GetURL("images/test.png")
	if url != "/uploads/images/test.png" {
		t.Errorf("GetURL = %s, want /uploads/images/test.png", url)
	}
}

func TestLocalStorage_Type(t *testing.T) {
	s := NewLocalStorage("/tmp", "/uploads")
	if s.Type() != "local" {
		t.Errorf("Type = %s, want local", s.Type())
	}
}

func TestDefaultStorage(t *testing.T) {
	tmpDir := t.TempDir()
	s := NewLocalStorage(tmpDir, "/uploads")
	Init(s)

	if Default == nil {
		t.Fatal("Default should not be nil after Init()")
	}
}
