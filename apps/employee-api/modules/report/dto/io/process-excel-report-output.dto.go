package io

import "github.com/google/uuid"

type ProcessExcelReportOutputDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Document string    `json:"document"`
}
