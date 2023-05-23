package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	ProjectID  string
	BucketName string
	UploadPath string
}

func InitializeSpaces() *storage.Client {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "secret.json") // FILL IN WITH YOUR FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("GCP Connected")

	return client

}

// func main() {
// 	//uploader.UploadFile("notes_test/abc.txt")
// 	r := gin.Default()
// 	r.POST("/upload", func(c *gin.Context) {
// 		f, err := c.FormFile("file_input")
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		blobFile, err := f.Open()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		err = uploader.UploadFile(blobFile, f.Filename)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "success",
// 		})
// 	})

// 	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }

// // UploadFile uploads an object
// func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
// 	ctx := context.Background()

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()

// 	// Upload an object with storage.Writer.
// 	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
// 	if _, err := io.Copy(wc, file); err != nil {
// 		return fmt.Errorf("io.Copy: %v", err)
// 	}
// 	if err := wc.Close(); err != nil {
// 		return fmt.Errorf("Writer.Close: %v", err)
// 	}

// 	return nil
// }
