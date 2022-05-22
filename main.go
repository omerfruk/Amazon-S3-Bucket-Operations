package main

import (
	"fmt"
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
	BASEDIR            string = "omerfruk-buclets"
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
	app.Get("/folders", ListFolder)

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

func ListFolder(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
