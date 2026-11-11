package transit

import (
	"context"
	"net/http"

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

	// S3 Client
	s3client := s3.NewFromConfig(awscfg,
		func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Upload.Endpoint)
		},
	)

	var transferMngr = &TransitManager{
		s3client:   s3client,
		httpclient: http.DefaultClient,
		cfg:        &cfg,
		ctx:        ctx,
	}

	return transferMngr, nil
}
