package storage

import (
	"mime/multipart"
)

type Storage interface {
	SaveImage(image multipart.File, imagePath string) error
	DeleteImage(imgPath string, thumbnailPath string) error
	ImageInit() error
}
