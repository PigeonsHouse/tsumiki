package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tsumiki/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MediaService interface {
	UploadAvatar(ctx context.Context, discordUserID string, imageURL string) (string, error)
	ResolveURL(path string) string
}

type mediaServiceImpl struct {
	s3Client  *s3.Client
	bucket    string
	publicURL string
}

func NewMediaService(ctx context.Context, s3Client *s3.Client) (MediaService, error) {
	svc := &mediaServiceImpl{
		s3Client:  s3Client,
		bucket:    env.S3Bucket,
		publicURL: env.S3PublicURL,
	}
	if err := svc.ensureBucket(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}

func (ms *mediaServiceImpl) ensureBucket(ctx context.Context) error {
	_, err := ms.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(ms.bucket),
	})
	if err != nil {
		var bae *types.BucketAlreadyExists
		var baoy *types.BucketAlreadyOwnedByYou
		if !errors.As(err, &bae) && !errors.As(err, &baoy) {
			return fmt.Errorf("バケットの作成に失敗しました: %w", err)
		}
	}
	return nil
}

func (ms *mediaServiceImpl) UploadAvatar(ctx context.Context, discordUserID string, imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("アバター画像の取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("アバター画像の読み込みに失敗しました: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}

	key := fmt.Sprintf("avatars/%s.png", discordUserID)

	_, err = ms.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ms.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("アバター画像のアップロードに失敗しました: %w", err)
	}

	return fmt.Sprintf("%s/%s", ms.bucket, key), nil
}

func (ms *mediaServiceImpl) ResolveURL(path string) string {
	return fmt.Sprintf("%s/%s", ms.publicURL, path)
}
