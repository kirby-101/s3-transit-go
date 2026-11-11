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
	ctx, cancel = context.WithCancel(context.Background())

	// Flags
	downloadHTTPURL = flag.String("url", "", "download source url")

	uploadBucket    = flag.String("r2-bucket", "", "R2 bucket name")
	uploadPath      = flag.String("r2-path", "", "R2 upload path")
	uploadEndpoint  = flag.String("r2-endpoint", "", "R2 upload endpoint")
	uploadAccessKey = flag.String("r2-acckey", "", "R2 Access-Key")
	uploadSecretKey = flag.String("r2-secretkey", "", "R2 Secret Key")
)

func main() {
	// parse argv
	flag.Parse()
	flag.VisitAll(
		func(_flag *flag.Flag) {
			if len(_flag.Value.String()) <= 0 {
				logError("[x] Argument: '-%s' cannot be empty\n", _flag.Name)
				os.Exit(1)
			} else {
				fmt.Printf("[cfg] %s = %s\n", _flag.Name, _flag.Value)
			}
		},
	)

	// bei SIGINT/SIGTERM cancel()
	go handleSignals(cancel)

	t, err := transit.New(
		ctx,
		transit.Config{
			URL: *downloadHTTPURL,
			Upload: transit.ConfigUpload{
				Bucket:    *uploadBucket,
				FilePath:  *uploadPath,
				Endpoint:  *uploadEndpoint,
				AccessKey: *uploadAccessKey,
				SecretKey: *uploadSecretKey,
			},
		},
	)
	if err != nil {
		logError("[x] TransitManager: %s\n", err.Error())
		os.Exit(1)
	}

	if err := t.Transfer(); err != nil {
		logError("[x] Transfer: %s\n", err.Error())
	}
}

func handleSignals(cancel context.CancelFunc) {
	var sigs = make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-sigs

	logError("[x] catched SIGINT/SIGTERM; aborting...")
	cancel()
}

func logError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
}
