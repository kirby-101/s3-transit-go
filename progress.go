package transit

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
)

type ProgressListener struct{}

func (p *ProgressListener) OnObjectTransferStart(_ context.Context, event *transfermanager.ObjectTransferStartEvent) {
	fmt.Printf("Transit Started: Total-Size: %s \n", humanBytes(uint64(event.TotalBytes)))
}

func (p *ProgressListener) OnObjectBytesTransferred(ctx context.Context, event *transfermanager.ObjectBytesTransferredEvent) {
	var (
		transferred = event.BytesTransferred
		total       = event.TotalBytes
	)

	if total > 0 {
		percent := float64(transferred) / float64(total) * 100
		fmt.Printf("\rProgress: %s / %s , (%.1f%%)", humanBytes(uint64(transferred)), humanBytes(uint64(total)), percent)

	} else {
		// unknown size
		fmt.Printf("\rProgress: %s", humanBytes(uint64(transferred)))
	}

	os.Stdout.Sync()
}

func (p *ProgressListener) OnObjectTransferComplete(_ context.Context, event *transfermanager.ObjectTransferCompleteEvent) {
	fmt.Println("Transit Succeeded.")
}

// retry?
func (p *ProgressListener) OnObjectTransferFailed(_ context.Context, event *transfermanager.ObjectTransferFailedEvent) {
	fmt.Printf("Transit failed: %v\n", event.Error.Error())
}
