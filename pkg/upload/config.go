package upload

import "github.com/caarlos0/env/v6"

// Configuration holds all of the options required to connect to minio or AWS S3.
type Configuration struct {
	MinIOEndpoint    string `env:"BUCKET_ENDPOINT"`
	MinIOAccessKeyID string `env:"BUCKET_ACCESS_KEY_ID"`
	MinIOAccessKey   string `env:"BUCKET_ACCESS_KEY"`
	UseSSL           bool   `env:"BUCKET_USE_SSL"`
	BucketName       string
	BucketLocation   string
}

// GetConfig parses configuration options from the environment, and returns a
// Configuration struct for further modification/amendment.
func GetConfig() *Configuration {
	environment := &Configuration{}
	if err := env.Parse(environment); err != nil {
		panic("unable to parse environment!")
	}

	return environment
}