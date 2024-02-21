package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func deleteFromS3(fileKey string, s3Client *s3.S3, bucket string, bucketUrl string) error {

	// URL에서 객체 키 추출
	key := extractKeyFromUrl(fileKey, bucket, bucketUrl)
	log.Println("key", fileKey)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// 에러 발생 시 처리 로직
	if err != nil {
		fmt.Printf("Failed to delete object from S3: %s, error: %v\n", fileKey, err)
	}

	return err
}

// URL에서 S3 객체 키를 추출하는 함수
func extractKeyFromUrl(url, bucket string, bucketUrl string) string {
	prefix := fmt.Sprintf("https://%s.%s/", bucket, bucketUrl)
	return strings.TrimPrefix(url, prefix)
}

func uploadImagesToS3(imgData []byte, thumbnailData []byte, contentType string, ext string, s3Client *s3.S3, bucket string, bucketUrl string, uid string) (string, string, error) {
	// 이미지 파일 이름과 썸네일 파일 이름 생성
	imgFileName := "images/diet/" + uid + "/images/" + uuid.New().String() + ext
	thumbnailFileName := "images/diet/" + uid + "/thumbnails/" + uuid.New().String() + ext

	// S3에 이미지 업로드
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imgFileName),
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", "", err
	}

	// S3에 썸네일 업로드
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(thumbnailFileName),
		Body:        bytes.NewReader(thumbnailData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", err
	}

	// 업로드된 이미지와 썸네일의 URL 생성 및 반환
	imgURL := "https://" + bucket + "." + bucketUrl + "/" + imgFileName
	thumbnailURL := "https://" + bucket + "." + bucketUrl + "/" + thumbnailFileName

	return imgURL, thumbnailURL, nil
}

func reduceImageSize(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	log.Println("image size: ", len(data))
	// 원본 이미지의 크기를 절반씩 줄이면서 10MB 이하로 만듦
	for len(data) > 10*1024*1024 {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		err := jpeg.Encode(&buf, resizedImg, nil)
		if err != nil {
			return nil, err
		}

		data = buf.Bytes()
		img = resizedImg
	}

	return data, nil
}

func createThumbnail(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 썸네일의 크기를 절반씩 줄이면서 1MB 이하로 만듦
	for {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		thumbnail := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, thumbnail, nil)
		if err != nil {
			return nil, err
		}

		thumbnailData := buf.Bytes()
		log.Println("thumbnailData size: ", len(thumbnailData))
		if len(thumbnailData) < 1024*1024 {
			return thumbnailData, nil
		}

		img = thumbnail
	}
}

func getImageFormat(imgData []byte) (contentType, extension string, err error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return "", "", err
	}

	contentType = "image/" + format
	extension = "." + format
	if format == "jpeg" {
		extension = ".jpg"
	}

	return contentType, extension, nil
}
