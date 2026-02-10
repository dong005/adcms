package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage MinIO/S3兼容存储（也可用于阿里云OSS S3兼容模式）
type MinIOStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
	baseURL  string // 外部访问URL前缀
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint  string // 如 "localhost:9000"
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	BaseURL   string // 外部访问URL，如 "https://cdn.example.com/bucket"
}

// NewMinIOStorage 创建MinIO存储实例
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("连接MinIO失败: %w", err)
	}

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查Bucket失败: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("创建Bucket失败: %w", err)
		}
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		scheme := "http"
		if cfg.UseSSL {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s/%s", scheme, cfg.Endpoint, cfg.Bucket)
	}

	return &MinIOStorage{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		baseURL:  strings.TrimRight(baseURL, "/"),
	}, nil
}

func (s *MinIOStorage) Type() string {
	return "minio"
}

func (s *MinIOStorage) Upload(file *multipart.FileHeader, dir string) (*FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	datePath := time.Now().Format("2006/01/02")
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	objectName := filepath.Join(dir, datePath, filename)
	objectName = strings.ReplaceAll(objectName, "\\", "/")

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx := context.Background()
	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("上传到MinIO失败: %w", err)
	}

	return &FileInfo{
		Name:     file.Filename,
		Path:     objectName,
		URL:      s.GetURL(objectName),
		Size:     file.Size,
		MimeType: contentType,
	}, nil
}

func (s *MinIOStorage) UploadReader(reader io.Reader, filename string, dir string) (*FileInfo, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	datePath := time.Now().Format("2006/01/02")
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	objectName := filepath.Join(dir, datePath, newFilename)
	objectName = strings.ReplaceAll(objectName, "\\", "/")

	ctx := context.Background()
	info, err := s.client.PutObject(ctx, s.bucket, objectName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("上传到MinIO失败: %w", err)
	}

	return &FileInfo{
		Name: filename,
		Path: objectName,
		URL:  s.GetURL(objectName),
		Size: info.Size,
	}, nil
}

func (s *MinIOStorage) Delete(path string) error {
	ctx := context.Background()
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("从MinIO删除失败: %w", err)
	}
	return nil
}

func (s *MinIOStorage) GetURL(path string) string {
	return s.baseURL + "/" + strings.ReplaceAll(path, "\\", "/")
}
