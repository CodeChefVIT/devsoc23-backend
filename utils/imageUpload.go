package utils

import (
	"context"
	"devsoc23-backend/helper"
	"devsoc23-backend/initializers"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
)

type PhotoForm struct {
	CampaignImage *multipart.FileHeader `form:"image" binding:"required"`
}

func UploadPhoto(payload *PhotoForm, S3Client *storage.Client) (string, error) {
	// Edit context
	ctx := context.TODO()

	BucketConfig := initializers.ClientUploader{
		ProjectID:  os.Getenv("PROJECT_ID"),
		BucketName: os.Getenv("BUCKET_NAME"),
		UploadPath: "images/",
	}

	file, err := payload.CampaignImage.Open()

	if err != nil {
		return "", fmt.Errorf("image could not be parsed: %w", err)
	}

	wc := S3Client.Bucket(BucketConfig.BucketName).Object(BucketConfig.UploadPath + helper.GenerateToken() + payload.CampaignImage.Filename).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		fmt.Println(err.Error())
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err.Error())
		return "", fmt.Errorf("Writer.Close: %v", err.Error())
	}
	u := ("https://storage.googleapis.com/" + BucketConfig.BucketName + "/" + wc.Attrs().Name)

	if err != nil {
		return "", fmt.Errorf("generating Image url failed: %w", err)
	}
	return u, nil

}
