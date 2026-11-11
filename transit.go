package transit

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
)

// 301?, User-Agent
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

	setDefaultHTTPHeaders(req)

	// Request
	resp, err := t.httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TransitManager) uploadToS3(bucket, filePath, contentType string, contentLength int64, src io.ReadCloser) error {
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
	if contentLength != 0 {
		uploadargs.ContentLength = aws.Int64(contentLength)
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
func (t *TransitManager) Transfer() error {
	var (
		contentType         = ""
		contentLength int64 = 0
	)

	/*
	 * Download
	 */

	// Get Download-Connection
	resp, err := t.downloadFromURL(t.cfg.URL)
	if err != nil {
		return fmt.Errorf("downloadFromURL: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("invalid response status-code: %d", resp.StatusCode)
	}

	/*
	* Upload
	 */

	// Content-Type
	contentType = resp.Header.Get("Content-Type")

	// Content-Lenght
	if resp.ContentLength != 0 {
		contentLength = resp.ContentLength
	}

	err = t.uploadToS3(
		t.cfg.Upload.Bucket,
		t.cfg.Upload.FilePath,
		contentType,
		contentLength,
		resp.Body,
	)

	return err
}
