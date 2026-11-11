package transit

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
)

// 301?
func (t *TransitManager) downloadFromURL(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		t.ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Request
	resp, err := t.httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TransitManager) uploadToS3(bucket, filePath, contentType string, contentLenght int64, src io.ReadCloser) error {
	var (
		uploadmanager = transfermanager.New(
			t.s3client,
			func(o *transfermanager.Options) {
				//o.PartSizeBytes = defaultChunkSize
				//o.Concurrency = defaultConcurrency
			},
		)
		uploadargs = &transfermanager.UploadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filePath),
			Body:   src,
		}
	)

	// Set Content-Lenght if available
	if contentLenght != 0 {
		uploadargs.ContentLength = aws.Int64(contentLenght)
	}

	// Set Content-Type if available
	if contentType != "" {
		uploadargs.ContentType = aws.String(contentType)
	}

	// S3 Upload
	_, err := uploadmanager.UploadObject(
		t.ctx,
		uploadargs,
		func(o *transfermanager.Options) {
			o.ObjectProgressListeners.Register(&ProgressListener{})
		},
	)

	return err
}

// ! logging,progress report
func (t *TransitManager) TransferFromURL(url, bucketName, filePath string) error {
	var (
		contentType         = ""
		contentLenght int64 = 0
	)

	/*
	 * Download
	 */

	// Get Download-Connection
	resp, err := t.downloadFromURL(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status-code: %d", resp.StatusCode)
	}

	/*
	* Upload
	 */

	// Content-Type
	contentType = resp.Header.Get("Content-Type")

	// Content-Lenght
	if resp.ContentLength != 0 {
		contentLenght = resp.ContentLength
	}

	return t.uploadToS3(
		bucketName,
		filePath,
		contentType,
		contentLenght,
		resp.Body,
	)
}
