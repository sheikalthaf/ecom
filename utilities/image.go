package utilities

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var s3Handler *S3Handler

func InitS3(bucket string) error {
	var err error
	s3Handler, err = NewS3Handler(bucket)
	return err
}

func UploadImage(c *fiber.Ctx, fieldName string, folderName string, oldImageName string) (*string, string, error) {
	imgOutput, err := processAndSaveImage(c, fieldName, folderName, nil)
	if err != nil {
		return nil, "", err
	}

	var oldImage *string
	newImage := oldImageName
	if imgOutput != nil {
		oldImage = &oldImageName
		newImage = imgOutput.CompressedImage
		*oldImage = ProperPathName(*oldImage)
	}
	return oldImage, ProperPathName(newImage), nil
}

func processAndSaveImage(c *fiber.Ctx, inputFileName string, outputFolderName string, outputFileName *string) (*ImageOutputViewModel, error) {
	file, err := c.FormFile(inputFileName)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, nil
	}
	if file.Size > 2000000 { // 2MB
		return nil, errors.New("file size is too large")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if outputFileName == nil {
		id := uuid.New().String()
		outputFileName = &id
	}

	replacer := strings.NewReplacer("\\", "/")

	output := &ImageOutputViewModel{
		CompressedImage: replacer.Replace(filepath.Join(outputFolderName, *outputFileName+".jpeg")),
		ThumbnailImage:  replacer.Replace(filepath.Join(outputFolderName, "thumbnails", *outputFileName+".jpeg")),
	}

	// Compress the image
	compressedImg, err := compressImage(src, 800)
	if err != nil {
		return nil, err
	}

	// Upload compressed image to S3
	err = s3Handler.UploadFile(context.Background(), output.CompressedImage, compressedImg)
	if err != nil {
		return nil, err
	}

	// Reset the file reader position
	src.Seek(0, 0)

	// Create a thumbnail
	thumbnailImg, err := createThumbnail(src, 150)
	if err != nil {
		return nil, err
	}

	// Upload thumbnail to S3
	err = s3Handler.UploadFile(context.Background(), output.ThumbnailImage, thumbnailImg)
	if err != nil {
		return nil, err
	}

	return output, nil
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

func DeleteImage(newImage string, oldImage *string, folder string, isError bool) {
	if oldImage != nil && isError {
		s3Handler.DeleteFile(context.Background(), newImage)
	} else if oldImage != nil && !isError {
		s3Handler.DeleteFile(context.Background(), *oldImage)
	}
}

// multipartFile is a helper struct to convert a bytes.Reader to a multipart.File
type multipartFile struct {
	*bytes.Reader
}

func (m *multipartFile) Close() error {
	return nil
}

type ImageOutputViewModel struct {
	CompressedImage string
	ThumbnailImage  string
}
