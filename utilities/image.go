package utilities

import (
	"bytes"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"ecom.com/app/storage"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	storage storage.Storage
}

// store the handler
var Image *Handler

func NewHandler(storage storage.Storage) {
	Image = &Handler{storage: storage}
}

func (h *Handler) UploadImage(c *fiber.Ctx, fieldName string, folderName string, oldImageName string) (*string, string, error) {
	// Open the file
	file, _ := c.FormFile(fieldName)
	// Check the file exist or not
	if file == nil {
		return nil, "", nil
	}
	if file.Size > 2000000 { // 2MB
		return nil, "", errors.New("file size is too large")
	}

	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}

	defer src.Close()

	// if outputFileName == nil {
	outputFileName := uuid.New().String()
	// outputFileName = &id
	// }

	output := &ImageOutputViewModel{
		CompressedImage: filepath.Join("images", folderName, outputFileName+".jpeg"),
		ThumbnailImage:  filepath.Join("images", folderName, "thumbnails", outputFileName+".jpeg"),
	}

	// Compress the image and save it to the "compressed" folder
	compressedImg, err := compressImage(src, 800)
	if err != nil {
		return nil, "", err
	}

	// Save the compressed image
	err = h.storage.SaveImage(compressedImg, output.CompressedImage)
	if err != nil {
		return nil, "", err
	}

	// Reset the file reader position
	src.Seek(0, 0)

	// Create a thumbnail
	thumbnailImg, err := createThumbnail(src, 150)
	if err != nil {
		return nil, "", err
	}
	// return output, nil
	// imgOutput, err := processAndSaveImage(c, fieldName, folderName, nil)
	// if err != nil {
	// 	return nil, "", err
	// }

	// Upload thumbnail to S3
	err = h.storage.SaveImage(thumbnailImg, output.ThumbnailImage)
	if err != nil {
		return nil, "", err
	}

	var oldImage *string
	newImage := oldImageName
	// if output != nil {
	oldImage = &oldImageName
	newImage = output.CompressedImage
	*oldImage = ProperPathName(*oldImage)
	// }
	return oldImage, ProperPathName(newImage), nil
}

func compressImage(img multipart.File, newWidth int) (multipart.File, error) {
	// Decode the image
	image, err := imaging.Decode(img)
	if err != nil {
		return nil, err
	}

	// Resize the image
	resizedImg := imaging.Resize(image, newWidth, 0, imaging.Lanczos)

	// Create a buffer to store the compressed image
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, resizedImg, imaging.JPEG, imaging.JPEGQuality(70))
	if err != nil {
		return nil, err
	}

	// Convert buffer to multipart.File
	return &multipartFile{bytes.NewReader(buf.Bytes())}, nil
}

func createThumbnail(img multipart.File, thumbnailSize int) (multipart.File, error) {
	// Decode the image
	image, err := imaging.Decode(img)
	if err != nil {
		return nil, err
	}

	// Create the thumbnail
	thumbnailImg := imaging.Thumbnail(image, thumbnailSize, thumbnailSize, imaging.Lanczos)

	// Create a buffer to store the thumbnail image
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, thumbnailImg, imaging.JPEG, imaging.JPEGQuality(90))
	if err != nil {
		return nil, err
	}

	// Convert buffer to multipart.File
	return &multipartFile{bytes.NewReader(buf.Bytes())}, nil
}

func (h *Handler) DeleteImage(newImage string, oldImage *string, folder string, isError bool) error {
	// delete the new images created
	var imgPath string
	var thumbnailPath string
	if oldImage != nil && isError {
		imgPath, thumbnailPath = deleteOldImages(newImage, folder)
	} else if oldImage != nil && !isError {
		// delete the old images
		imgPath, thumbnailPath = deleteOldImages(*oldImage, folder)
	}

	return h.storage.DeleteImage(imgPath, thumbnailPath)
}

func deleteOldImages(imgPath string, folderName string) (string, string) {
	// Check if the image path is empty or not ends with the default image
	if imgPath == "" || !isImageFileName(imgPath) {
		return "", ""
	}
	// replace folderName with folderName + thumbnails
	thumbnailPath := strings.Replace(imgPath, folderName, folderName+"/thumbnails", 1)
	// os.Remove(imgPath)
	// os.Remove(thumbnailPath)
	return imgPath, thumbnailPath
}

func isImageFileName(fileName string) bool {
	// List of valid image file extensions
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	// Get the file extension
	ext := strings.ToLower(filepath.Ext(fileName))

	// Check if the file extension is in the list of valid image file extensions
	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}

	return false
}

type ImageOutputViewModel struct {
	CompressedImage string
	ThumbnailImage  string
}

// multipartFile is a helper struct to convert a bytes.Reader to a multipart.File
type multipartFile struct {
	*bytes.Reader
}

func (m *multipartFile) Close() error {
	return nil
}
