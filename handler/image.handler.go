package handlers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadImage(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "server error"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "No file is received"})
	}

	file := files[0]
	fileSize := file.Size
	if fileSize < 10*1024 || fileSize > 20*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"message": "File size must be between 10 KB and 20 MB"})
	}
	fmt.Println(fileSize)

	fileType := file.Header.Get("Content-Type")
	fmt.Println(fileType)
	if fileType != "image/jpeg" && fileType != "image/jpg" {
		return c.Status(400).JSON(fiber.Map{"error": "File must be a JPEG image"})
	}

	fileName := uuid.New().String() + filepath.Ext(file.Filename)

	S3_REGION := os.Getenv("AWS_REGION")
	S3_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	S3_SECRET_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3_BUCKET_NAME := os.Getenv("AWS_S3_BUCKET_NAME")
	// fmt.Println("AWS_REGION:", S3_REGION)
	// fmt.Println("AWS_ACCESS_KEY:", S3_ID)
	// fmt.Println("AWS_SECRET_ACCESS_KEY:", S3_SECRET_KEY)
	// fmt.Println("AWS_S3_BUCKET_NAME:", S3_BUCKET_NAME)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(S3_ID, S3_SECRET_KEY, ""),
	})
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "Failed to create AWS session"})
	}

	s3Client := s3.New(sess)

	fileBuffer, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to open file"})
	}
	defer fileBuffer.Close()

	// Copy the fileBuffer to buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, fileBuffer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to read file"})
	}

	// Decode the image from the buffer
	_, _, err = image.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to decode image"})
	}

	// Upload the file to S3.
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(S3_BUCKET_NAME),
		Key:    aws.String(fileName),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(buffer.Bytes()),  // Re-read from the buffer
	})

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to upload file to S3"})
	}
	url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", S3_BUCKET_NAME, S3_REGION, fileName)

	return c.Status(200).JSON(fiber.Map{"message": "success", "data": fiber.Map{"imageUrl": url}})
}
