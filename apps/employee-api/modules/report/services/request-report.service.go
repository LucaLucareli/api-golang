package reports

import (
	"context"

	"employee-api/modules/report/dto/io"
	"shared/queue/enqueue"
	"shared/queue/payloads"
	"shared/types"

	"github.com/hibiken/asynq"
)

type RequestReportService struct {
	client *asynq.Client
}

func NewRequestReportService(client *asynq.Client) *RequestReportService {
	return &RequestReportService{
		client: client,
	}
}

func (s *RequestReportService) Execute(
	ctx context.Context,
	input io.RequestReportInputDTO,
) error {

	payload := payloads.GenerateReportPayload{
		ReportID:   1,
		UserID:     input.UserID,
		ReportType: types.ReportType(input.ReportTypeID),
	}

	return enqueue.EnqueueGenerateReport(s.client, payload)
}
