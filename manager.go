package transit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TransitManager struct {
	s3client   *s3.Client
	httpclient *http.Client
	cfg        *Config
	ctx        context.Context
}

func New(ctx context.Context, cfg Config) (*TransitManager, error) {
	// S3 Config
	awscfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.Upload.AccessKey,
				cfg.Upload.SecretKey,
				"",
			),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	var (
		// S3 Client
		s3client = s3.NewFromConfig(awscfg,
			func(o *s3.Options) {
				o.BaseEndpoint = aws.String(cfg.Upload.Endpoint)
			},
		)
		// HTTP Client
		httpclient = &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				IdleConnTimeout:       10 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// maximum of 10 redirects
				if len(via) >= 10 {
					return fmt.Errorf(">10 redirects. can't continue")
				}

				//? log
				fmt.Printf("Followed Redirect")

				// headers
				for k, v := range via[0].Header {
					if _, exists := req.Header[k]; !exists {
						req.Header[k] = v
					}
				}

				return nil
			},
			Timeout: 0,
		}
	)

	var transferMngr = &TransitManager{
		s3client:   s3client,
		httpclient: httpclient,
		cfg:        &cfg,
		ctx:        ctx,
	}

	return transferMngr, nil
}
