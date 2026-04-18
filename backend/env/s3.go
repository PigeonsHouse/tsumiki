package env

import (
	"fmt"
	"os"
)

var (
	S3Endpoint        string
	S3PublicURL       string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string
)

func LoadS3Env() error {
	S3Endpoint = os.Getenv("S3_ENDPOINT")
	if S3Endpoint == "" {
		return fmt.Errorf("loading env error: S3_ENDPOINT")
	}
	S3PublicURL = os.Getenv("S3_PUBLIC_URL")
	if S3PublicURL == "" {
		return fmt.Errorf("loading env error: S3_PUBLIC_URL")
	}
	S3Bucket = os.Getenv("S3_BUCKET")
	if S3Bucket == "" {
		return fmt.Errorf("loading env error: S3_BUCKET")
	}
	S3AccessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
	if S3AccessKeyID == "" {
		return fmt.Errorf("loading env error: S3_ACCESS_KEY_ID")
	}
	S3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	if S3SecretAccessKey == "" {
		return fmt.Errorf("loading env error: S3_SECRET_ACCESS_KEY")
	}

	return nil
}
