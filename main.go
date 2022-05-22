package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	REGION             string = "eu-central-1"
	S3CREDENTIALID     string = "AKIAYU3HHMJTEBFC2IZC"
	S3CREDENTIALSECRET string = "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"
)

var Sess = session.Session{}

func main() {
	app := fiber.New()

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(REGION),
			Credentials: credentials.NewStaticCredentials(
				S3CREDENTIALID,
				S3CREDENTIALSECRET,
				"",
			),
		})
	if err != nil {
		panic(err)
	}

	Sess = *sess

	app.Get("/buckets", ListBuckets)
	app.Get("/create-bucket", CreateBuckets)

	app.Get("/list-bucket-items", ListBucketsItems)
	app.Post("/upload-file-to-bucket", UploadFileToBucket)
	app.Get("/download-file-to-bucket", DownloadFileToBucket)
	app.Get("/delete-file-from-bucket", DeleteItemInBucket)
	app.Get("/delete-bucket", DeleteBucket)

	app.Listen(":3000")
}

func ListBuckets(c *fiber.Ctx) error {
	svc := s3.New(&Sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return c.JSON(result.Buckets)
}

func CreateBuckets(c *fiber.Ctx) error {
	bucket := c.Query("bucket-name")

	svc := s3.New(&Sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket:                    aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String(REGION)},
	})
	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)

	return c.JSON("Bucket created successfuly bucket name: " + bucket)
}

func ListBucketsItems(c *fiber.Ctx) error {
	bucket := c.Query("bucket-name")

	// Create S3 service client
	svc := s3.New(&Sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")

	return c.JSON(resp.Contents)
}

type FileDis struct {
	FileName string `json:"file_name"`
	Bucket   string `json:"bucket"`
}

func UploadFileToBucket(c *fiber.Ctx) error {
	var fileDis FileDis
	err := c.BodyParser(&fileDis)
	if err != nil {
		exitErrorf("body parse hatasi ")
	}

	headerFile, err := c.FormFile(fileDis.FileName)
	if err != nil {
		exitErrorf("dosya okunamadı")
	}

	file, err := headerFile.Open()
	uploader := s3manager.NewUploader(&Sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fileDis.Bucket),
		Key:    aws.String(fileDis.FileName),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", fileDis.FileName, fileDis.Bucket, err)
	}
	fmt.Printf("Successfully uploaded %q to %q\n", fileDis.FileName, fileDis.Bucket)

	return c.JSON(fileDis)
}

func DownloadFileToBucket(c *fiber.Ctx) error {

	bucket := c.Query("bucket")
	item := c.Query("item")

	file, err := os.Create(item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(&Sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return c.Download(file.Name())
}

func DeleteItemInBucket(c *fiber.Ctx) error {
	bucket := c.Query("bucket")
	item := c.Query("item")

	// Create S3 service client
	svc := s3.New(&Sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(item)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", item, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	fmt.Printf("Object %q successfully deleted\n", item)
	return c.JSON("Deleting success")
}

func DeleteBucket(c *fiber.Ctx) error {
	bucket := c.Query("bucket")
	// Create S3 service client
	svc := s3.New(&Sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Unable to delete bucket %q, %v", bucket, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucket)

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be deleted, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully deleted\n", bucket)
	return c.JSON("Bucket deleting success")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
