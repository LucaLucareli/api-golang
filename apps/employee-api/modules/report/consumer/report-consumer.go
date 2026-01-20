package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"shared/types"

	"shared/queue/payloads"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type ReportProcessorFunc func(
	payload payloads.GenerateReportPayload,
) error

type ReportProcessorRegistry struct {
	processors map[types.ReportType]ReportProcessorFunc
}

func NewReportProcessorRegistry(
	ExcelProcessor ReportProcessorFunc,
) *ReportProcessorRegistry {
	return &ReportProcessorRegistry{
		processors: map[types.ReportType]ReportProcessorFunc{
			types.Excel: ExcelProcessor,
		},
	}
}

func (r *ReportProcessorRegistry) Process(
	payload payloads.GenerateReportPayload,
) error {
	processor, ok := r.processors[payload.ReportType]
	if !ok {
		return fmt.Errorf(
			"Nenhum processador registrado para o tipo de relatório %d",
			payload.ReportType,
		)
	}

	return processor(payload)
}

type ReportProcessor struct {
	registry *ReportProcessorRegistry
}

func NewReportProcessor(
	registry *ReportProcessorRegistry,
) *ReportProcessor {
	return &ReportProcessor{
		registry: registry,
	}
}

func (p *ReportProcessor) Handler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		return p.ProcessTask(ctx, t)
	}
}

func (p *ReportProcessor) ProcessTask(
	ctx context.Context,
	t *asynq.Task,
) error {
	var payload payloads.GenerateReportPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Info().
		Int32("report_id", payload.ReportID).
		Str("user_id", payload.UserID).
		Int("report_type", int(payload.ReportType)).
		Msg("Iniciando geração de relatório")

	if err := p.registry.Process(payload); err != nil {
		log.Error().
			Err(err).
			Int32("report_id", payload.ReportID).
			Msg("Falha ao gerar relatório")

		return err
	}

	log.Info().
		Int32("report_id", payload.ReportID).
		Msg("Relatório gerado com sucesso")

	return nil
}
