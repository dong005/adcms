package storage

import (
	"io"
	"mime/multipart"
)

// FileInfo 上传后的文件信息
type FileInfo struct {
	Name     string `json:"name"`      // 原始文件名
	Path     string `json:"path"`      // 存储路径（相对）
	URL      string `json:"url"`       // 访问URL
	Size     int64  `json:"size"`      // 文件大小
	MimeType string `json:"mime_type"` // MIME类型
}

// Storage 文件存储接口
type Storage interface {
	// Upload 上传文件
	Upload(file *multipart.FileHeader, dir string) (*FileInfo, error)
	// UploadReader 从 Reader 上传
	UploadReader(reader io.Reader, filename string, dir string) (*FileInfo, error)
	// Delete 删除文件
	Delete(path string) error
	// GetURL 获取文件访问URL
	GetURL(path string) string
	// Type 存储类型标识
	Type() string
}

// 全局存储实例
var Default Storage

// Init 初始化默认存储
func Init(s Storage) {
	Default = s
}
