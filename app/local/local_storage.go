package local

import (
	"io"
	"mime/multipart"
	"os"
)

type LocalImageStorage struct{}

func NewLocalImageStorage() *LocalImageStorage {
	return &LocalImageStorage{}
}

func (l *LocalImageStorage) AppendUrl(imagePath string) string {
	return imagePath
}

func (l *LocalImageStorage) SaveImage(c multipart.File, imageName string) error {
	// Implement local file system upload logic
	dst, err := os.Create(imageName)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, c); err != nil {
		return err
	}
	return nil
}

func (l *LocalImageStorage) DeleteImage(imgPath string, thumbnailPath string) error {
	os.Remove(imgPath)
	os.Remove(thumbnailPath)
	return nil
}

func (l *LocalImageStorage) ImageInit() error {
	os.MkdirAll("images/product/thumbnails/", 0755)
	return nil
}
