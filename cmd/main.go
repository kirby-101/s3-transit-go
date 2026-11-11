package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	transit "github.com/kirby-101/s3-transit-go"
)

var (
	// Flags
	downloadHTTPURL = *flag.String("url", "", "download source url")

	uploadBucket    = *flag.String("r2-bucket", "", "R2 bucket name")
	uploadPath      = *flag.String("r2-path", "", "R2 upload path")
	uploadEndpoint  = *flag.String("r2-endpoint", "", "R2 upload endpoint")
	uploadAccountID = *flag.String("r2-accid", "", "R2 Account-ID")
	uploadAccessKey = *flag.String("r2-acckey", "", "R2 Access-Key")
	uploadSecretKey = *flag.String("r2-secretkey", "", "R2 Secret Key")
)

func init() {
	flag.Parse() // parse argv

	flag.VisitAll(
		func(_flag *flag.Flag) {
			if len(_flag.Value.String()) <= 0 {
				fmt.Fprintf(os.Stderr, "Argument: '-%s' cannot be empty\n", _flag.Name)
				os.Exit(1)
			}
		},
	)
}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())

	// bei SIGINT/SIGTERM cancel()
	go handleSignals(cancel)

	transit.New(
		ctx,
		transit.Config{
			URL: downloadHTTPURL,
			Upload: transit.ConfigUpload{
				Bucket:    uploadBucket,
				FilePath:  uploadPath,
				Endpoint:  uploadEndpoint,
				AccountID: uploadAccountID,
				AccessKey: uploadAccessKey,
				SecretKey: uploadSecretKey,
			},
		},
	)
}

func handleSignals(cancel context.CancelFunc) {
	var sigs = make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-sigs
	logError("catched SIGINT/SIGTERM")

	cancel()
}

func logError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
}
