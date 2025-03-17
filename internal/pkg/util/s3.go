package util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"os"
	"sync"
)

type Storage struct {
	bucket   string
	domain   string
	access   string
	secret   string
	sess     *session.Session
	svc      *s3.S3
	uploader *s3manager.Uploader
	once     sync.Once
}

func NewStorage(bucket, domain, access, secret string) *Storage {
	return &Storage{
		bucket: bucket,
		domain: domain,
		access: access,
		secret: secret,
	}
}

func (s *Storage) initSession() error {
	var err error
	s.once.Do(func() {
		s.sess, err = session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials(s.access, s.secret, ""),
			Region:           aws.String("default"),
			Endpoint:         aws.String(s.domain),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err != nil {
			log.Printf("خطا در ایجاد سشن: %v", err)
			return
		}
		s.svc = s3.New(s.sess)
		s.uploader = s3manager.NewUploader(s.sess)
	})
	return err
}

func (s *Storage) ListFiles() error {
	if err := s.initSession(); err != nil {
		return err
	}

	resp, err := s.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("error in list files: %v", err)
	}

	for _, item := range resp.Contents {
		fmt.Printf("%s (%d bytes)\n", *item.Key, *item.Size)
	}

	fmt.Printf("📌 number of files : %d\n", len(resp.Contents))
	return nil
}

func (s *Storage) UploadFile(file io.Reader, destPath string) error {

	if err := s.initSession(); err != nil {
		return err
	}

	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(destPath),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("error in upload file %s: %v", destPath, err)
	}

	fmt.Printf("file %s , uploaded in path %s", destPath, destPath)
	return nil
}

func (s *Storage) DownloadFile(remotePath, localPath string) error {
	if err := s.initSession(); err != nil {
		return err
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("خطا در ایجاد فایل %s: %v", localPath, err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(s.sess)
	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(remotePath),
	})
	if err != nil {
		return fmt.Errorf("خطا در دانلود فایل %s: %v", remotePath, err)
	}

	fmt.Printf("✅ دانلود شد: %s (%d bytes)\n", localPath, numBytes)
	return nil
}

func (s *Storage) DeleteFile(remotePath string) error {
	if err := s.initSession(); err != nil {
		return err
	}

	_, err := s.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(remotePath),
	})
	if err != nil {
		return fmt.Errorf("خطا در حذف فایل %s: %v", remotePath, err)
	}

	fmt.Printf("🗑️ فایل %s حذف شد.\n", remotePath)
	return nil
}
