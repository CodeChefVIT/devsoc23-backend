package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type PhotoForm struct {
	CampaignImage *multipart.FileHeader `form:"mage" binding:"required"`
}

func readFile(file *multipart.FileHeader) ([]byte, error) {
	openedFile, _ := file.Open()

	binaryFile, err := io.ReadAll(openedFile)

	if err != nil {
		return nil, err
	}

	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			log.Fatalf("Failed closing file %v", file.Filename)
		}
	}(openedFile)
	return binaryFile, nil
}

func UploadPhoto(payload *PhotoForm, s3ClientLoaded *s3.S3) (string, error) {
	bucketName := os.Getenv("DO_SPACE_NAME")

	binaryImageFile, err := readFile(payload.CampaignImage)

	if err != nil {
		return "", fmt.Errorf("generating Binary Image failed: %w", err)
	}

	tImageUuid := uuid.New()

	imagePath := "/raisze/" + tImageUuid.String() + "-" + payload.CampaignImage.Filename

	tfbImage := bytes.NewReader(binaryImageFile)
	tImageobject := s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imagePath),
		Body:   tfbImage,
		ACL:    aws.String("public-read"),
	}
	_, uploadErr := s3ClientLoaded.PutObject(&tImageobject)
	if uploadErr != nil {
		return "", fmt.Errorf("failed image upload: %w", err)

	}

	campaignImageUrl := "https://spaces-shortsqueeze.sgp1.digitaloceanspaces.com" + imagePath

	if err != nil {
		return "", fmt.Errorf("generating Image url failed: %w", err)
	}
	return campaignImageUrl, nil

}
