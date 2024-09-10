package storage

import (
	"mime/multipart"
)

type Storage interface {
	SaveImage(image multipart.File, imagePath string) error
	AppendUrl(imagePath string) string
	DeleteImage(imgPath string, thumbnailPath string) error
	ImageInit() error
}
