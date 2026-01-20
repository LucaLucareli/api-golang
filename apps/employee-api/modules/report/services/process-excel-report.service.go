package reports

import (
	"context"
	"fmt"
	"shared"
	"shared/queue/payloads"
	"shared/report"

	"github.com/rs/zerolog/log"
)

func NewProcessExcelReport(state shared.AppState) func(payload payloads.GenerateReportPayload) error {
	return func(payload payloads.GenerateReportPayload) error {
		ctx := context.Background()
		return ProcessExcelReport(ctx, payload, state)
	}
}

func ProcessExcelReport(
	ctx context.Context,
	payload payloads.GenerateReportPayload,
	state shared.AppState,
) error {
	log.Info().Int32("report_id", payload.ReportID).Msg("Processando relatório EXCEL")

	path := fmt.Sprintf("report_%d.xlsx", payload.ReportID)

	writer, err := report.NewExcelWriter(path)
	if err != nil {
		log.Error().Err(err).Msg("Falha ao criar ExcelWriter")
		return err
	}

	headers := []string{"ID", "Nome", "Documento"}
	if err := writer.WriteHeader(headers); err != nil {
		log.Error().Err(err).Msg("Falha ao escrever cabeçalho")
		return err
	}

	rowsStream, err := state.UserRepo.FindManyUsersToReport(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Falha ao iniciar streaming de usuários")
		return err
	}

	for row := range rowsStream {
		if row.Err != nil {
			log.Error().Err(row.Err).Msg("Erro ao ler usuário do stream")
			continue
		}

		values := []any{
			row.User.ID,
			row.User.Name,
			row.User.Document,
		}

		if err := writer.WriteRow(values); err != nil {
			log.Error().Err(err).Msg("Falha ao escrever linha")
			return err
		}
	}

	if err := writer.Close(); err != nil {
		log.Error().Err(err).Msg("Falha ao salvar arquivo Excel")
		return err
	}

	log.Info().Msgf("Relatório Excel gerado com sucesso: %s", path)
	return nil
}
