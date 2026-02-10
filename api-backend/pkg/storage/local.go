package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorage 本地文件存储
type LocalStorage struct {
	BasePath string // 本地存储根目录，如 "./uploads"
	BaseURL  string // 访问URL前缀，如 "/uploads" 或 "http://cdn.example.com/uploads"
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  strings.TrimRight(baseURL, "/"),
	}
}

func (s *LocalStorage) Type() string {
	return "local"
}

func (s *LocalStorage) Upload(file *multipart.FileHeader, dir string) (*FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	datePath := time.Now().Format("2006/01/02")
	relDir := filepath.Join(dir, datePath)
	absDir := filepath.Join(s.BasePath, relDir)

	if err := os.MkdirAll(absDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relPath := filepath.Join(relDir, filename)
	absPath := filepath.Join(s.BasePath, relPath)

	dst, err := os.Create(absPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	return &FileInfo{
		Name:     file.Filename,
		Path:     relPath,
		URL:      s.GetURL(relPath),
		Size:     file.Size,
		MimeType: file.Header.Get("Content-Type"),
	}, nil
}

func (s *LocalStorage) UploadReader(reader io.Reader, filename string, dir string) (*FileInfo, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	datePath := time.Now().Format("2006/01/02")
	relDir := filepath.Join(dir, datePath)
	absDir := filepath.Join(s.BasePath, relDir)

	if err := os.MkdirAll(absDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relPath := filepath.Join(relDir, newFilename)
	absPath := filepath.Join(s.BasePath, relPath)

	dst, err := os.Create(absPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, reader)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	return &FileInfo{
		Name: filename,
		Path: relPath,
		URL:  s.GetURL(relPath),
		Size: written,
	}, nil
}

func (s *LocalStorage) Delete(path string) error {
	absPath := filepath.Join(s.BasePath, path)
	if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetURL(path string) string {
	return s.BaseURL + "/" + strings.ReplaceAll(path, "\\", "/")
}
