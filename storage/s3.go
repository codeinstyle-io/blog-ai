package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

// S3Provider implements Provider interface for AWS S3 storage
type S3Provider struct {
	client *s3.Client
	bucket string
}

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	// TODO: Logo endpoint

	// fallback to default
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// NewS3Provider creates a new S3Provider
func NewS3Provider(bucket, region, endpoint, access_key, secret_key string) (*S3Provider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access_key, secret_key, "")),
		config.WithRegion(region),
	)

	// Create custom resolver for endpoint if provided

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.EndpointResolverV2 = &resolverV2{}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	return &S3Provider{
		client: client,
		bucket: bucket,
	}, nil
}

// Save implements Provider.Save
func (p *S3Provider) Save(filename string, reader io.Reader) (string, error) {
	// Generate unique filename with slugified name
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	filename = fmt.Sprintf("%d-%s%s", time.Now().Unix(), slugify(name), ext)

	// Upload to S3
	_, err := p.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			fmt.Printf("Failed to upload file %s to S3: %s (%s)\n", filename, ae.ErrorMessage(), ae.ErrorCode())
		} else {
			fmt.Printf("Failed to upload file %s to S3: %v\n", filename, err)
		}
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded file %s to S3 bucket %s\n", filename, p.bucket)
	return filename, nil
}

// Delete implements Provider.Delete
func (p *S3Provider) Delete(path string) error {
	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			fmt.Printf("Failed to delete file %s from S3: %s (%s)\n", path, ae.ErrorMessage(), ae.ErrorCode())
		} else {
			fmt.Printf("Failed to delete file %s from S3: %v\n", path, err)
		}
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	fmt.Printf("Successfully deleted file %s from S3 bucket %s\n", path, p.bucket)
	return nil
}

// Get implements Provider.Get
func (p *S3Provider) Get(path string) (io.ReadCloser, error) {
	result, err := p.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			fmt.Printf("Failed to get file %s from S3: %s (%s)\n", path, ae.ErrorMessage(), ae.ErrorCode())
			if ae.ErrorCode() == "NoSuchKey" {
				return nil, fmt.Errorf("file not found: %s", path)
			}
		} else {
			fmt.Printf("Failed to get file %s from S3: %v\n", path, err)
		}
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	return result.Body, nil
}
