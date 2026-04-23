package media

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"tsumiki/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MediaService interface {
	UploadAvatar(ctx context.Context, userID int, r io.Reader, contentType string) (string, error)
	UploadTsumikiMedia(ctx context.Context, tsumikiID int, r io.Reader, contentType string, ext string) (string, error)
	UploadThumbnail(ctx context.Context, userID int, r io.Reader, contentType string, ext string) (string, error)
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

func (ms *mediaServiceImpl) UploadAvatar(ctx context.Context, userID int, r io.Reader, contentType string) (string, error) {
	if contentType == "" {
		contentType = "image/png"
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("アバター画像の読み込みに失敗しました: %w", err)
	}

	hash := sha256.Sum256(data)
	key := fmt.Sprintf("avatars/%d/%x.png", userID, hash)

	_, err = ms.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ms.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("アバター画像のアップロードに失敗しました: %w", err)
	}

	return key, nil
}

func (ms *mediaServiceImpl) UploadTsumikiMedia(ctx context.Context, tsumikiID int, r io.Reader, contentType string, ext string) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	hash := sha256.Sum256(data)
	key := fmt.Sprintf("tsumikis/%d/medias/%x%s", tsumikiID, hash, ext)

	_, err = ms.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ms.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %w", err)
	}

	return key, nil
}

func (ms *mediaServiceImpl) UploadThumbnail(ctx context.Context, userID int, r io.Reader, contentType string, ext string) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("サムネイルの読み込みに失敗しました: %w", err)
	}
	hash := sha256.Sum256(data)
	key := fmt.Sprintf("thumbnails/%d/%x%s", userID, hash, ext)
	_, err = ms.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ms.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("サムネイルのアップロードに失敗しました: %w", err)
	}
	return key, nil
}

func (ms *mediaServiceImpl) ResolveURL(path string) string {
	return fmt.Sprintf("%s/%s/%s", ms.publicURL, ms.bucket, path)
}
