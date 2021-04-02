package forklift

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Uploader uploads arbitrary byte strings to S3
type S3Uploader struct {
	templateURI string
	pipeMap     map[string]*io.PipeWriter
	uploader    *s3manager.Uploader
	wg          sync.WaitGroup
}

// AddRecord determines where each record should get written to
// And then sends a Write command
func (s3u *S3Uploader) AddRecord(r Record) {
	w := s3u.getWriter(r)
	w.Write(r.GetBytes())
}

func (s3u *S3Uploader) newOutputFile(path string) *io.PipeWriter {
	r1, w1 := io.Pipe()
	s3u.wg.Add(1)
	go func() {
		// Upload the file to S3.
		log.Printf("Creating new S3 output file: %s", path)
		bucket, key := getS3BucketAndKey(path)
		result, err := s3u.uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   r1,
		})
		if err != nil {
			fmt.Printf("failed to upload file, %v", err)
		}
		log.Printf("Completing file upload to, %s", result.Location)
		s3u.wg.Done()
	}()

	return w1
}

func getS3BucketAndKey(path string) (string, string) {
	u, _ := url.Parse(path)
	return u.Host, strings.TrimLeft(u.Path, "/")
}

func (s3u *S3Uploader) getWriter(r Record) *io.PipeWriter {
	outputkey := r.FormatPath(s3u.templateURI)
	if val, ok := s3u.pipeMap[outputkey]; ok {
		return val
	}

	s3u.pipeMap[outputkey] = s3u.newOutputFile(outputkey)
	return s3u.pipeMap[outputkey]
}

// Close waits for all data to finish uploading
func (s3u *S3Uploader) Close() {
	// First close all of our writers
	for _, w := range s3u.pipeMap {
		w.Close()
	}

	// Then wait for them to finish!
	s3u.wg.Wait()
}

func bucketFinder(b string) string {
	svc := s3.New(session.New(&aws.Config{
		Region:                         aws.String("us-east-1"),
		DisableRestProtocolURICleaning: aws.Bool(true),
	}))

	resp, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: &b,
	})
	if err != nil {
		panic(err)
	}

	// If this does *not* return a region, then it is in us-east-1
	region := "us-east-1"
	if resp.LocationConstraint != nil {
		region = *resp.LocationConstraint
	}

	return region
}

// NewS3Uploader initializes the S3 Uploader and returns a pointer to S3Uploader
// in order to satisfy the interface requirements.
func NewS3Uploader(u string) *S3Uploader {
	log.Printf("Connecting to S3 with template: %s", u)

	bucket, _ := getS3BucketAndKey(u)
	region := bucketFinder((bucket))
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))

	return &S3Uploader{
		templateURI: u,
		uploader:    s3manager.NewUploader(sess),
		wg:          sync.WaitGroup{},
		pipeMap:     map[string]*io.PipeWriter{},
	}
}
