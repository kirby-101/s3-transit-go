package transit

import (
	"context"
	"fmt"

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
		fmt.Printf("Progress: %s / %s , (%.1f%%)", humanBytes(uint64(transferred)), humanBytes(uint64(total)), percent)

	} else {
		// unknown size
		fmt.Printf("Progress: %s", humanBytes(uint64(transferred)))

	}
}

func (p *ProgressListener) OnObjectTransferComplete(_ context.Context, event *transfermanager.ObjectTransferCompleteEvent) {
	fmt.Println("Transit Succeeded.")
}

// retry?
func (p *ProgressListener) OnObjectTransferFailed(_ context.Context, event *transfermanager.ObjectTransferFailedEvent) {
	fmt.Printf("Transit failed: %v\n", event.Error.Error())
}

func humanBytes(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
