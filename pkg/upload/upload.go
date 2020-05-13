package upload

import (
	"io"
	"os"

	"github.com/minio/minio-go/v6"
)

var (
	client     *minio.Client
	bucketName string
)

// Configuration holds all of the options required to connect to minio or AWS S3.
type Configuration struct {
	MinIOEndpoint    string
	MinIOAccessKeyID string
	MinIOAccessKey   string
	UseSSL           bool
	BucketName       string
	BucketLocation   string
}

// Connect makes a connection to the object storage destination, and ensures that the
// specified bucket is available - if it doesn't exist, it shall be created.
func Connect(conf *Configuration) error {
	c, err := minio.New(conf.MinIOEndpoint, conf.MinIOAccessKeyID, conf.MinIOAccessKey, conf.UseSSL)
	if err != nil {
		return err
	}

	client = c
	bucketName = conf.BucketName

	hasBucket, err := client.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if !hasBucket {
		if err := client.MakeBucket(bucketName, conf.BucketLocation); err != nil {
			return err
		}
	}

	return nil
}

// Upload is a wrapper around the underlying storage library, and uploads a file
// from an io.Reader, returning the number of bytes uploaded and any errors.
func Upload(file io.Reader, fileInfo os.FileInfo) (int64, error) {
	return client.PutObject(bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
}
