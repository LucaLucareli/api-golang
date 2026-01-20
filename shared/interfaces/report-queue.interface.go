package interfaces

import (
	"context"
	"shared/queue/payloads"
)

type ReportQueue interface {
	Enqueue(ctx context.Context, payload payloads.GenerateReportPayload) error
}
